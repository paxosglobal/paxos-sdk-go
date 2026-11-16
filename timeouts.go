package paxos

import (
	"context"
	"fmt"
	"time"
)

// WithDefaultCallTimeout bounds each logical call, including authentication,
// retries, backoff and the full response body. The SDK default is 30 seconds.
func WithDefaultCallTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d <= 0 {
			c.optionErr = fmt.Errorf("call timeout must be positive; use WithoutDefaultCallTimeout to disable")
			return
		}
		c.callTimeout = d
	}
}

// WithoutDefaultCallTimeout removes the SDK call deadline. Context cancellation still applies.
func WithoutDefaultCallTimeout() Option { return func(c *Client) { c.callTimeout = 0 } }

// WithDefaultCallAttemptTimeout bounds each HTTP attempt, including OAuth attempts.
// By default no additional attempt cap is applied.
func WithDefaultCallAttemptTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d <= 0 {
			c.optionErr = fmt.Errorf("call attempt timeout must be positive")
			return
		}
		c.callAttemptTimeout = d
	}
}

// WithoutDefaultCallAttemptTimeout removes the additional attempt cap.
func WithoutDefaultCallAttemptTimeout() Option { return func(c *Client) { c.callAttemptTimeout = 0 } }

// WithCallTimeout overrides the total call budget, including authentication and retries.
func WithCallTimeout(d time.Duration) RequestOption {
	return func(o *requestOptions) {
		if d <= 0 {
			o.optionErr = fmt.Errorf("call timeout must be positive; use WithoutCallTimeout to disable")
		}
		o.callTimeout = &d
	}
}

// WithoutCallTimeout removes the SDK total deadline for this call.
func WithoutCallTimeout() RequestOption {
	return func(o *requestOptions) { d := time.Duration(0); o.callTimeout = &d }
}

// WithCallAttemptTimeout overrides the per HTTP attempt budget for this call.
func WithCallAttemptTimeout(d time.Duration) RequestOption {
	return func(o *requestOptions) {
		if d <= 0 {
			o.optionErr = fmt.Errorf("call attempt timeout must be positive")
		}
		o.callAttemptTimeout = &d
	}
}

// WithoutCallAttemptTimeout removes the additional attempt cap for this call.
func WithoutCallAttemptTimeout() RequestOption {
	return func(o *requestOptions) { d := time.Duration(0); o.callAttemptTimeout = &d }
}

func waitBackoff(ctx context.Context, attempt int) error {
	timer := time.NewTimer(calculateBackoff(attempt))
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
