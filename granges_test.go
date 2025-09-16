package granges_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AyakuraYuki/granges"
)

func TestOpen(t *testing.T) {
	r := granges.Open(4, 8)
	assert.False(t, r.IsInvalid())
	checkContains(t, r)
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 4, r.LowerEndpoint())
	assert.EqualValues(t, granges.OPEN, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 8, r.UpperEndpoint())
	assert.EqualValues(t, granges.OPEN, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "(4..8)", r.String())
}

func TestOpen_invalid(t *testing.T) {
	r := granges.Open(4, 1)
	assert.True(t, r.IsInvalid())
	r, err := granges.OpenE(3, 3)
	assert.Error(t, err)
	assert.True(t, r.IsInvalid())
}

func TestClosed(t *testing.T) {
	r := granges.Closed(5, 7)
	assert.False(t, r.IsInvalid())
	checkContains(t, r)
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 5, r.LowerEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 7, r.UpperEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "[5..7]", r.String())
}

func TestClosed_invalid(t *testing.T) {
	r, err := granges.ClosedE(4, 3)
	assert.Error(t, err)
	assert.True(t, r.IsInvalid())
}

func TestClosedOpen(t *testing.T) {
	r := granges.ClosedOpen(5, 8)
	assert.False(t, r.IsInvalid())
	checkContains(t, r)
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 5, r.LowerEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 8, r.UpperEndpoint())
	assert.EqualValues(t, granges.OPEN, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "[5..8)", r.String())
}

func TestOpenClosed(t *testing.T) {
	r := granges.OpenClosed(4, 7)
	assert.False(t, r.IsInvalid())
	checkContains(t, r)
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 4, r.LowerEndpoint())
	assert.EqualValues(t, granges.OPEN, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 7, r.UpperEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "(4..7]", r.String())
}

func checkContains(t *testing.T, r granges.Range[int]) {
	assert.False(t, r.Contains(4))
	assert.True(t, r.Contains(5))
	assert.True(t, r.Contains(7))
	assert.False(t, r.Contains(8))
}

func TestSingleton(t *testing.T) {
	r := granges.Singleton(4)
	assert.False(t, r.Contains(3))
	assert.True(t, r.Contains(4))
	assert.False(t, r.Contains(5))
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 4, r.LowerEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 4, r.UpperEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "[4..4]", r.String())
}

func TestLessThan(t *testing.T) {
	r := granges.LessThan(5)
	assert.True(t, r.Contains(math.MinInt))
	assert.True(t, r.Contains(4))
	assert.False(t, r.Contains(5))
	assertUnboundedBelow(t, r)
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 5, r.UpperEndpoint())
	assert.EqualValues(t, granges.OPEN, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "(-\u221e..5)", r.String())
}

func TestGreaterThan(t *testing.T) {
	r := granges.GreaterThan(5)
	assert.False(t, r.Contains(5))
	assert.True(t, r.Contains(6))
	assert.True(t, r.Contains(math.MaxInt))
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 5, r.LowerEndpoint())
	assert.EqualValues(t, granges.OPEN, r.LowerBoundType())
	assertUnboundedAbove(t, r)
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "(5..+\u221e)", r.String())
}

func TestAtLeast(t *testing.T) {
	r := granges.AtLeast(6)
	assert.False(t, r.Contains(5))
	assert.True(t, r.Contains(6))
	assert.True(t, r.Contains(math.MaxInt))
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 6, r.LowerEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.LowerBoundType())
	assertUnboundedAbove(t, r)
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "[6..+\u221e)", r.String())
}

func TestAtMost(t *testing.T) {
	r := granges.AtMost(4)
	assert.True(t, r.Contains(math.MinInt))
	assert.True(t, r.Contains(4))
	assert.False(t, r.Contains(5))
	assertUnboundedBelow(t, r)
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 4, r.UpperEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.UpperBoundType())
	assert.False(t, r.IsEmpty())
	assert.EqualValues(t, "(-\u221e..4]", r.String())
}

func TestAll(t *testing.T) {
	all := granges.All[int]()
	assert.True(t, all.Contains(math.MinInt))
	assert.True(t, all.Contains(math.MaxInt))
	assertUnboundedBelow(t, all)
	assertUnboundedAbove(t, all)
	assert.False(t, all.IsEmpty())
	assert.EqualValues(t, "(-\u221e..+\u221e)", all.String())
}

func assertUnboundedBelow(t *testing.T, r granges.Range[int]) {
	assert.False(t, r.HasLowerBound())
	_, err := r.LowerEndpointE()
	assert.ErrorIs(t, err, granges.ErrRangeSideUnbounded)
	_, err = r.LowerBoundTypeE()
	assert.ErrorIs(t, err, granges.ErrUnboundedCut)
}

func assertUnboundedAbove(t *testing.T, r granges.Range[int]) {
	assert.False(t, r.HasUpperBound())
	_, err := r.UpperEndpointE()
	assert.ErrorIs(t, err, granges.ErrRangeSideUnbounded)
	_, err = r.UpperBoundTypeE()
	assert.ErrorIs(t, err, granges.ErrUnboundedCut)
}
