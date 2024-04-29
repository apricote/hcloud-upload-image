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
