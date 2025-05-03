// SPDX-License-Identifier: MPL-2.0
// From https://github.com/hetznercloud/terraform-provider-hcloud/blob/v1.46.1/internal/control/retry.go
// Copyright (c) Hetzner Cloud GmbH

package control

import (
	"context"
	"time"

	"github.com/apricote/hcloud-upload-image/hcloudimages/backoff"
	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
)

// Retry executes f at most maxTries times.
func Retry(ctx context.Context, maxTries int, f func() error) error {
	logger := contextlogger.From(ctx)

	var err error

	backoffFunc := backoff.ExponentialBackoffWithLimit(2, 200*time.Millisecond, 2*time.Second)

	for try := 0; try < maxTries; try++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = f()
		if err != nil {
			sleep := backoffFunc(try)
			logger.DebugContext(ctx, "operation failed, waiting before trying again", "try", try, "backoff", sleep)
			time.Sleep(sleep)
			continue
		}

		return nil
	}

	return err
}
