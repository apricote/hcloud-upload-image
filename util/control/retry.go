package control

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// From https://github.com/hetznercloud/terraform-provider-hcloud/blob/v1.46.1/internal/control/retry.go

// ExponentialBackoffWithLimit returns a [hcloud.BackoffFunc] which implements an exponential
// backoff.
// It uses the formula:
//
//	min(b^retries * d, limit)
func ExponentialBackoffWithLimit(b float64, d time.Duration, limit time.Duration) hcloud.BackoffFunc {
	return func(retries int) time.Duration {
		current := time.Duration(math.Pow(b, float64(retries))) * d

		if current > limit {
			return limit
		} else {
			return current
		}
	}
}

// DefaultRetries is a constant for the maximum number of retries we usually do.
// However, callers of Retry are free to choose a different number.
const DefaultRetries = 5

type abortErr struct {
	Err error
}

func (e abortErr) Error() string {
	return e.Err.Error()
}

func (e abortErr) Unwrap() error {
	return e.Err
}

// AbortRetry aborts any further attempts of retrying an operation.
//
// If err is passed Retry returns the passed error. If nil is passed, Retry
// returns nil.
func AbortRetry(err error) error {
	if err == nil {
		return nil
	}
	return abortErr{Err: err}
}

// Retry executes f at most maxTries times.
func Retry(ctx context.Context, maxTries int, f func() error) error {
	var err error

	backoff := ExponentialBackoffWithLimit(2, 1*time.Second, 30*time.Second)

	for try := 0; try < maxTries; try++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		var aerr abortErr

		err = f()
		if errors.As(err, &aerr) {
			return aerr.Err
		}
		if err != nil {
			sleep := backoff(try)
			time.Sleep(sleep)
			continue
		}

		return nil
	}

	return err
}
