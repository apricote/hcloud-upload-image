// SPDX-License-Identifier: MPL-2.0
// From https://github.com/hetznercloud/terraform-provider-hcloud/blob/v1.46.1/internal/control/retry.go
// Copyright (c) Hetzner Cloud GmbH

package backoff

import (
	"math"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

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
