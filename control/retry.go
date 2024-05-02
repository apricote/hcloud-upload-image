package control

import (
	"context"
	"math"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/apricote/hcloud-upload-image/contextlogger"
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

// Retry executes f at most maxTries times.
func Retry(ctx context.Context, maxTries int, f func() error) error {
	logger := contextlogger.From(ctx)

	var err error

	backoff := ExponentialBackoffWithLimit(2, 1*time.Second, 30*time.Second)

	for try := 0; try < maxTries; try++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = f()
		if err != nil {
			sleep := backoff(try)
			logger.DebugContext(ctx, "operation failed, waiting before trying again", "try", try, "backoff", sleep)
			time.Sleep(sleep)
			continue
		}

		return nil
	}

	return err
}
