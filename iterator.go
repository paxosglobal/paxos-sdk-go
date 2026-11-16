package paxos

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"time"
)

type iterMode int

const (
	modeNone iterMode = iota
	modeItems
	modePages
)

var errMixedIteration = errors.New("iterator: item iteration and Pages cannot be mixed on one iterator")

// pageFetcher is a function that fetches a page of results given a cursor.
// It returns the items, the next cursor (empty string if no more pages), and an error.
type pageFetcher[T any] func(ctx context.Context, cursor string) (items []T, nextCursor string, err error)

// Iterator provides generic pagination support using Go generics.
// It lazily fetches pages one at a time to avoid memory accumulation.
//
// Termination is guaranteed even if the server misbehaves: iteration stops
// on an empty page (even if a cursor was still returned) or a cursor
// identical to the one just requested, in addition to the normal case of an
// absent cursor.
//
// Iterator is NOT safe for concurrent use. Each goroutine should use its own iterator.
type Iterator[T any] struct {
	ctx     context.Context
	cancel  context.CancelFunc
	fetcher pageFetcher[T]

	items      []T
	itemIndex  int
	nextCursor string

	maxDuration  time.Duration
	started      bool
	maxItems     int
	itemsYielded int

	pageMeta func() ResponseMetadata
	mode     iterMode

	err    error
	done   bool
	closed bool
}

// iteratorOptions holds optional bounds on an Iterator's total output.
type iteratorOptions struct {
	err         error
	maxItems    int
	maxDuration time.Duration
}

// iteratorOption configures optional bounds when constructing an Iterator
// directly via newIterator. Service methods generated from the OpenAPI spec
// take RequestOption instead (WithMaxItems, WithMaxDuration) and translate
// them into IteratorOptions internally.
type iteratorOption func(*iteratorOptions)

// withIteratorMaxItems caps the total number of items the Iterator will
// yield, regardless of how many pages the fetcher has left to offer.
func withIteratorMaxItems(n int) iteratorOption {
	return func(o *iteratorOptions) {
		o.maxItems = n
		if n <= 0 {
			o.err = fmt.Errorf("max items must be positive")
		}
	}
}

// withIteratorMaxDuration bounds the total wall-clock time the Iterator may
// spend across all of its fetcher calls. Once exceeded, Next returns false
// and Err reports an error wrapping context.DeadlineExceeded.
func withIteratorMaxDuration(d time.Duration) iteratorOption {
	return func(o *iteratorOptions) {
		o.maxDuration = d
		if d <= 0 {
			o.err = fmt.Errorf("max duration must be positive")
		}
	}
}

// newIterator creates a new iterator that uses the provided fetcher function
// to lazily load pages of results.
func newIterator[T any](ctx context.Context, fetcher pageFetcher[T], opts ...iteratorOption) *Iterator[T] {
	var o iteratorOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &Iterator[T]{
		ctx:         ctx,
		err:         o.err,
		maxDuration: o.maxDuration,
		fetcher:     fetcher,
		nextCursor:  "", // Empty cursor means start from the beginning
		maxItems:    o.maxItems,
	}
}

// stop releases the iterator's internal deadline timer, if any.
// Safe to call more than once.
func (it *Iterator[T]) stop() {
	if it.cancel != nil {
		it.cancel()
		it.cancel = nil
	}
}

// Close releases resources held by the iterator — specifically, the
// internal timer backing a WithMaxDuration deadline, if one was set.
//
// Iteration releases the timer automatically once it completes (whether by
// exhaustion, error, or maxItems), so Close is only needed when a caller stops
// consuming an iterator early before it naturally finishes. Close also marks
// the iterator closed: Next returns false and Pages yields nothing afterwards.
// Safe to call multiple times, and safe to call even if no WithMaxDuration was
// set.
func (it *Iterator[T]) Close() {
	it.closed = true
	it.release()
}

func (it *Iterator[T]) release() {
	it.done = true
	it.stop()
}

// begin runs the checks shared by Next and Pages before any further work and
// reports whether iteration may continue.
func (it *Iterator[T]) begin() bool {
	if it.err != nil || it.done {
		return false
	}

	if !it.started {
		it.started = true
		if it.maxDuration > 0 {
			it.ctx, it.cancel = context.WithTimeout(it.ctx, it.maxDuration)
		}
	}
	if err := it.ctx.Err(); err != nil {
		it.err = withRequestID(err, "")
		it.stop()
		return false
	}

	// Enforce maxItems before doing any more work — including before
	// fetching a page we no longer need, if the cap happens to land exactly
	// on a page boundary.
	if it.maxItems > 0 && it.itemsYielded >= it.maxItems {
		it.done = true
		it.stop()
		return false
	}
	return true
}

// fetchPage fetches the next page, applying the termination guarantees shared
// by Next and Pages. It returns false when iteration failed; an empty page is
// returned with true after marking the iterator done.
func (it *Iterator[T]) fetchPage() ([]T, bool) {
	sentCursor := it.nextCursor
	var items []T
	items, it.nextCursor, it.err = it.fetcher(it.ctx, sentCursor)
	if it.err != nil {
		it.stop()
		return nil, false
	}

	if err := it.ctx.Err(); err != nil {
		it.err = withRequestID(err, "")
		it.stop()
		return nil, false
	}

	// Guard against a non-advancing cursor: if the server echoes back
	// exactly the cursor we just sent, following it again would loop
	// forever. Take this page's items, but don't fetch any further.
	if it.nextCursor != "" && it.nextCursor == sentCursor {
		it.nextCursor = ""
	}

	// An empty page terminates iteration even if it still carries a
	// next_page_cursor, instead of looping on empty pages forever. The page
	// itself is returned so Pages can surface it.
	if len(items) == 0 {
		it.done = true
		it.stop()
		return items, true
	}
	return items, true
}

// Next advances the iterator to the next item.
// It returns true if there is an item available, false if the iteration is complete or an error occurred.
// Call Err() after Next returns false to check for errors.
func (it *Iterator[T]) Next() bool {
	if it.mode == modePages && it.err == nil && !it.closed {
		it.err = errMixedIteration
		it.stop()
		return false
	}
	if !it.begin() {
		return false
	}
	it.mode = modeItems

	// If we have items remaining in the current page, advance to next item
	if it.itemIndex < len(it.items) {
		it.itemIndex++
		if it.itemIndex < len(it.items) {
			it.itemsYielded++
			return true
		}
	}

	// Current page exhausted, check if there are more pages
	if it.items != nil && it.nextCursor == "" {
		// We've already fetched at least one page and there's no next cursor
		it.done = true
		it.stop()
		return false
	}

	items, ok := it.fetchPage()
	if !ok || len(items) == 0 {
		return false
	}

	// Store the new page and reset index
	it.items = items
	it.itemIndex = 0
	it.itemsYielded++
	return true
}

// Page is one fetched page of results.
type Page[T any] struct {
	// Items are the page's items, trimmed to the remaining WithMaxItems budget.
	Items []T

	// NextCursor is the cursor for the following page. It is empty on the last
	// page, except that an empty page reports the cursor as received. It is
	// also empty on a page trimmed by WithMaxItems (resuming would skip the
	// trimmed items). A page stopped by the repeated-cursor guard has an empty
	// NextCursor, which takes precedence over the empty-page rule when the
	// cursor equals the one just sent. A page that ends exactly on the
	// WithMaxItems cap keeps its cursor.
	NextCursor string

	// Metadata is this page's own response metadata, independent of any
	// WithResponseMetadata holder.
	Metadata ResponseMetadata
}

// Pages returns a lazy range-over-func sequence of pages. A page is fetched
// only when the consumer advances to it, and the sequence shares the
// iterator's termination guarantees (empty page, repeated cursor, WithMaxItems,
// WithMaxDuration, context cancellation). An empty page, including an empty
// first page or one that still carries a cursor, is yielded once with its own
// Metadata and NextCursor as received (empty if it repeats the cursor just
// sent, including the starting PageCursor), and iteration then stops. A fetch
// failure is yielded as the final pair with a zero Page.
//
// Breaking out of the loop does not close the iterator: calling Pages again
// continues from the next unfetched page, and yields nothing once exhausted.
// Call Close if you will not resume and WithMaxDuration was set. To resume in
// a later list call, pass the last yielded NextCursor as PageCursor.
//
// Item iteration (Next, Current, Collect, All) and Pages cannot be mixed on one
// live iterator. Calling Pages after item iteration has started yields a single
// (zero Page, error) pair without poisoning the iterator, so Next still works.
// Next, Collect or All after Pages has started, even once it has finished
// normally, returns false (or the error) and sets Err. A failed or closed
// iterator instead ends in either view: Next returns false with Err holding the
// original error (nil after a plain Close), and Pages yields nothing once
// closed, or the stored error once per range after a failure, with no request
// and no mixing error (TypeScript and Python end silently there; the error
// stays on Err).
//
// A WithResponseMetadata holder still reflects only the final page fetched.
func (it *Iterator[T]) Pages() iter.Seq2[Page[T], error] {
	return func(yield func(Page[T], error) bool) {
		if it.closed {
			return
		}
		if it.mode == modeItems {
			if it.err != nil {
				yield(Page[T]{}, it.err)
			} else {
				yield(Page[T]{}, errMixedIteration)
			}
			return
		}
		it.mode = modePages
		for {
			if !it.begin() {
				break
			}
			if it.items != nil && it.nextCursor == "" {
				it.done = true
				it.stop()
				break
			}
			items, ok := it.fetchPage()
			if !ok {
				break
			}
			if n := len(items); n > 0 {
				if it.maxItems > 0 && it.itemsYielded+n > it.maxItems {
					items = items[:it.maxItems-it.itemsYielded]
					it.nextCursor = ""
				}
				it.items = items
				it.itemIndex = len(items) - 1
				it.itemsYielded += len(items)
			}
			page := Page[T]{Items: items, NextCursor: it.nextCursor}
			if it.pageMeta != nil {
				page.Metadata = it.pageMeta()
			}
			if !yield(page, nil) {
				return
			}
		}
		if err := it.Err(); err != nil {
			yield(Page[T]{}, err)
		}
	}
}

// Current returns the current item.
// Only valid after Next() returns true.
func (it *Iterator[T]) Current() T {
	if it.itemIndex >= 0 && it.itemIndex < len(it.items) {
		return it.items[it.itemIndex]
	}
	var zero T
	return zero
}

// Err returns the error that stopped the iteration, if any. If a
// WithMaxDuration deadline was exceeded, this wraps context.DeadlineExceeded
// (check with errors.Is).
func (it *Iterator[T]) Err() error {
	return it.err
}

// Collect drains and releases the iterator. If iteration fails it returns the
// items collected before the failure together with the error. A normally
// exhausted iterator is not closed, so Pages still reports the mixing error.
func (it *Iterator[T]) Collect() ([]T, error) {
	defer it.release()
	var out []T
	for it.Next() {
		out = append(out, it.Current())
	}
	return out, it.Err()
}

// All returns a range-over-func sequence whose final pair carries any
// iteration error; items fetched before the failure are yielded first. It
// releases the iterator when it finishes and closes it when the loop breaks.
func (it *Iterator[T]) All() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		defer it.release()
		for it.Next() {
			if !yield(it.Current(), nil) {
				it.Close()
				return
			}
		}
		if err := it.Err(); err != nil {
			var zero T
			yield(zero, err)
		}
	}
}

// iteratorOptionsFrom translates the maxItems/maxDuration fields of a parsed
// requestOptions into IteratorOptions, for the pagination entry points
// (paginateTyped, Client.paginate) that accept RequestOption rather than
// iteratorOption directly.
func iteratorOptionsFrom(reqOpts requestOptions) []iteratorOption {
	opts := make([]iteratorOption, 0, 3)
	if reqOpts.optionErr != nil {
		opts = append(opts, func(o *iteratorOptions) { o.err = reqOpts.optionErr })
	}
	if reqOpts.maxItems > 0 {
		opts = append(opts, withIteratorMaxItems(reqOpts.maxItems))
	}
	if reqOpts.maxDuration > 0 {
		opts = append(opts, withIteratorMaxDuration(reqOpts.maxDuration))
	}
	return opts
}
