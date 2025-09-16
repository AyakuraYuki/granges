package granges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkContains(t *testing.T, r Range[int]) {
	assert.False(t, r.Contains(4))
	assert.True(t, r.Contains(5))
	assert.True(t, r.Contains(7))
	assert.False(t, r.Contains(8))
}
