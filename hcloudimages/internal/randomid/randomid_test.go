// SPDX-License-Identifier: MIT
// From https://gitlab.com/hetznercloud/fleeting-plugin-hetzner/-/blob/0f60204582289c243599f8ca0f5be4822789131d/internal/utils/random_test.go
// Copyright (c) 2024 Hetzner Cloud GmbH

package randomid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomID(t *testing.T) {
	found1, err := Generate()
	assert.NoError(t, err)
	found2, err := Generate()
	assert.NoError(t, err)

	assert.Len(t, found1, 8)
	assert.Len(t, found2, 8)
	assert.NotEqual(t, found1, found2)
}
