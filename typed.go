package paxos

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// getTyped performs a GET request and deserializes the response into a typed struct.
func getTyped[T any](c *Client, ctx context.Context, path string, opts ...RequestOption) (*T, error) {
	var result T
	_, err := c.do(ctx, http.MethodGet, path, nil, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// getWithParamsTyped performs a GET request with query parameters and deserializes
// the response into a typed struct. The params argument can be a
// map[string]any or a pointer to a typed params struct.
func getWithParamsTyped[T any](c *Client, ctx context.Context, path string, params any, opts ...RequestOption) (*T, error) {
	path = appendQueryParams(path, toParamsMap(params))
	var result T
	_, err := c.do(ctx, http.MethodGet, path, nil, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// postTyped performs a POST request and deserializes the response into a typed struct.
func postTyped[T any](c *Client, ctx context.Context, path string, body any, opts ...RequestOption) (*T, error) {
	var result T
	_, err := c.do(ctx, http.MethodPost, path, body, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// putTyped performs a PUT request and deserializes the response into a typed struct.
func putTyped[T any](c *Client, ctx context.Context, path string, body any, opts ...RequestOption) (*T, error) {
	var result T
	_, err := c.do(ctx, http.MethodPut, path, body, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// deleteTyped performs a DELETE request and deserializes the response into a typed struct.
func deleteTyped[T any](c *Client, ctx context.Context, path string, opts ...RequestOption) (*T, error) {
	var result T
	_, err := c.do(ctx, http.MethodDelete, path, nil, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// deleteWithParamsTyped performs a DELETE request with query parameters and deserializes
// the response into a typed struct. The params argument can be a
// map[string]any or a pointer to a typed params struct.
func deleteWithParamsTyped[T any](c *Client, ctx context.Context, path string, params any, opts ...RequestOption) (*T, error) {
	path = appendQueryParams(path, toParamsMap(params))
	var result T
	_, err := c.do(ctx, http.MethodDelete, path, nil, &result, opts...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// paginateTyped returns an iterator that auto-paginates and yields typed items.
// It deserializes items directly from the raw JSON response, avoiding an
// intermediate map[string]any round-trip. The itemsField parameter specifies
// the JSON field name that contains the array of items (e.g. "payments").
// The params argument can be a map[string]any or a pointer to a typed params struct.
// Pass WithMaxItems and/or WithMaxDuration to bound the total iteration.
func paginateTyped[T any](c *Client, ctx context.Context, path string, itemsField string, params any, opts ...RequestOption) *Iterator[T] {
	reqOpts := parseRequestOptions(opts...)

	baseParams := toParamsMap(params)
	callOpts, pageMeta := pageMetadataCapture(reqOpts, opts)

	fetcher := func(ctx context.Context, cursor string) ([]T, string, error) {
		// Copy params so we don't mutate the caller's map.
		p := make(map[string]any)
		for k, v := range baseParams {
			p[k] = v
		}
		if cursor != "" {
			p["page_cursor"] = cursor
		} else {
			delete(p, "page_cursor")
		}

		// Use json.RawMessage so c.do gives us the raw JSON bytes
		// instead of decoding into map[string]any first.
		var raw json.RawMessage
		fullPath := appendQueryParams(path, p)
		meta, err := c.do(ctx, http.MethodGet, fullPath, nil, &raw, callOpts...)
		if err != nil {
			return nil, "", err
		}

		// Parse the top-level object to extract next_page_cursor and the
		// items array, without fully decoding every nested value.
		var envelope map[string]json.RawMessage
		if err := json.Unmarshal(raw, &envelope); err != nil {
			return nil, "", newMalformedResponseError(meta.StatusCode, meta.RequestID, err)
		}

		// Extract next_page_cursor from the envelope.
		var nextCursor string
		if cursorRaw, ok := envelope["next_page_cursor"]; ok {
			// Ignore errors; cursor may be null or absent.
			_ = json.Unmarshal(cursorRaw, &nextCursor)
		}

		// Extract items from the specified field.
		var items []T
		if fieldRaw, ok := envelope[itemsField]; ok {
			var err error
			if items, err = decodeItems[T](fieldRaw); err != nil {
				return nil, "", newMalformedResponseError(meta.StatusCode, meta.RequestID, err)
			}
		}

		return items, nextCursor, nil
	}

	it := newIterator(ctx, fetcher, iteratorOptionsFrom(reqOpts)...)
	it.nextCursor = startCursor(baseParams)
	it.pageMeta = pageMeta
	return it
}

// startCursor returns the caller's page_cursor, so the repeated-cursor guard
// also compares against the cursor the first request is sent with.
func startCursor(params map[string]any) string {
	cursor, _ := params["page_cursor"].(string)
	return cursor
}

// pageMetadataCapture returns the options for one page fetch and a function
// that snapshots that page's metadata. The caller's holder, if any, is read
// rather than replaced so its final-page semantics are unchanged.
func pageMetadataCapture(reqOpts requestOptions, opts []RequestOption) ([]RequestOption, func() ResponseMetadata) {
	holder := reqOpts.responseMetadata
	callOpts := opts
	if holder == nil {
		holder = &ResponseMetadata{}
		callOpts = append(opts[:len(opts):len(opts)], WithResponseMetadata(holder))
	}
	return callOpts, func() ResponseMetadata {
		m := *holder
		m.RawBody = bytes.Clone(holder.RawBody)
		return m
	}
}

// toParamsMap converts a params argument to map[string]any.
// Accepts nil, map[string]any, or a pointer to a struct with form/json tags.
// Nil pointers and empty slices are omitted. All other set fields are included
// (including zero-value scalars like 0 or false when the field is not a pointer).
// Panics if params is an unsupported type to prevent silently dropping query parameters.
func toParamsMap(params any) map[string]any {
	if params == nil {
		return nil
	}

	// Fast path: already a map.
	if m, ok := params.(map[string]any); ok {
		return m
	}

	// Struct pointer path: use reflection to extract set fields.
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("toParamsMap: unsupported params type %T; expected nil, map[string]any, or struct pointer", params))
	}

	m := make(map[string]any)
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		fv := v.Field(i)

		// Use the form tag as the canonical key (query parameter name).
		// Fall back to json tag if form is not set.
		tag := field.Tag.Get("form")
		if tag == "" || tag == "-" {
			tag = field.Tag.Get("json")
		}
		if tag == "" || tag == "-" {
			continue
		}
		key := strings.Split(tag, ",")[0]

		// Skip nil pointers (optional fields not set).
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		// Skip empty slices.
		if fv.Kind() == reflect.Slice && fv.Len() == 0 {
			continue
		}

		// Normalize the value for query string serialization.
		m[key] = normalizeParamValue(fv)
	}

	return m
}

// normalizeParamValue converts typed values to primitives suitable for
// query string serialization via appendQueryParams.
func normalizeParamValue(rv reflect.Value) any {
	val := rv.Interface()

	if t, ok := val.(time.Time); ok {
		return t.UTC().Format(time.RFC3339Nano)
	}

	// Check TextMarshaler first (handles time.Time and similar).
	if tm, ok := val.(encoding.TextMarshaler); ok {
		if text, err := tm.MarshalText(); err == nil {
			return string(text)
		}
	}

	// Handle slices by converting elements to []string for appendQueryParams.
	// Named-string types (e.g. []CryptoNetwork) have Kind()==String but are not
	// []string, so appendQueryParams would format them as "[BITCOIN ETHEREUM]".
	if rv.Kind() == reflect.Slice {
		elemType := rv.Type().Elem()
		if elemType.Kind() != reflect.Interface && rv.Type() != reflect.TypeOf([]string{}) {
			result := make([]string, rv.Len())
			for i := range rv.Len() {
				elem := rv.Index(i)
				if t, ok := elem.Interface().(time.Time); ok {
					result[i] = t.UTC().Format(time.RFC3339Nano)
					continue
				}
				if tm, ok := elem.Interface().(encoding.TextMarshaler); ok {
					if text, err := tm.MarshalText(); err == nil {
						result[i] = string(text)
						continue
					}
				}
				result[i] = fmt.Sprintf("%v", elem.Interface())
			}
			return result
		}
	}

	// Check Stringer (handles typed enums like SortOrder).
	if s, ok := val.(fmt.Stringer); ok {
		return s.String()
	}

	return val
}
