package granges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut_IsLessThan(t *testing.T) {
	belowAll := NewBelowAll[int]()
	assert.True(t, belowAll.IsLessThan(2))

	aboveAll := NewAboveAll[int]()
	assert.False(t, aboveAll.IsLessThan(2))

	belowValue := NewBelowValue[int](1)
	assert.True(t, belowValue.IsLessThan(2))
	assert.True(t, belowValue.IsLessThan(1))
	assert.False(t, belowValue.IsLessThan(0))

	aboveValue := NewAboveValue[int](1)
	assert.True(t, aboveValue.IsLessThan(2))
	assert.False(t, aboveValue.IsLessThan(1))
	assert.False(t, aboveValue.IsLessThan(0))

	// -1 means any invalid cut type
	// surely this will not happen on the caller side
	invalidCut := Cut[int]{cutType: -1}
	assert.False(t, invalidCut.IsLessThan(0))
}

func TestCut_DescribeAsLowerBound(t *testing.T) {
	belowAll := NewBelowAll[int]()
	assert.EqualValues(t, "(-∞", belowAll.DescribeAsLowerBound())

	aboveAll := NewAboveAll[int]()
	assert.EqualValues(t, "", aboveAll.DescribeAsLowerBound())

	belowValue := NewBelowValue[int](1)
	assert.EqualValues(t, "[1", belowValue.DescribeAsLowerBound())

	aboveValue := NewAboveValue[int](1)
	assert.EqualValues(t, "(1", aboveValue.DescribeAsLowerBound())

	// -1 means any invalid cut type
	// surely this will not happen on the caller side
	invalidCut := Cut[int]{cutType: -1}
	assert.EqualValues(t, "", invalidCut.DescribeAsLowerBound())
}

func TestCut_DescribeAsUpperBound(t *testing.T) {
	belowAll := NewBelowAll[int]()
	assert.EqualValues(t, "", belowAll.DescribeAsUpperBound())

	aboveAll := NewAboveAll[int]()
	assert.EqualValues(t, "+∞)", aboveAll.DescribeAsUpperBound())

	belowValue := NewBelowValue[int](1)
	assert.EqualValues(t, "1)", belowValue.DescribeAsUpperBound())

	aboveValue := NewAboveValue[int](1)
	assert.EqualValues(t, "1]", aboveValue.DescribeAsUpperBound())

	// -1 means any invalid cut type
	// surely this will not happen on the caller side
	invalidCut := Cut[int]{cutType: -1}
	assert.EqualValues(t, "", invalidCut.DescribeAsUpperBound())
}

func TestOrderingCuts(t *testing.T) {
	a := LessThan(0).lowerBound
	b := AtLeast(0).lowerBound
	c := GreaterThan(0).lowerBound
	d := AtLeast(1).lowerBound
	e := GreaterThan(1).lowerBound
	f := GreaterThan(1).upperBound

	testCompareToAndEquals(t, []Cut[int]{a, b, c, d, e, f})
}

func testCompareToAndEquals(t *testing.T, rs []Cut[int]) {
	for i := range rs {
		v := rs[i]
		for j := 0; j < i; j++ {
			lesser := rs[j]
			assert.True(t, lesser.Compare(v) < 0)
			assert.False(t, lesser.Compare(v) == 0)
		}

		assert.EqualValues(t, 0, v.Compare(v)) // self compare
		assert.True(t, v.Compare(v) == 0)

		for j := i + 1; j < len(rs); j++ {
			greater := rs[j]
			assert.True(t, greater.Compare(v) > 0)
			assert.False(t, greater.Compare(v) == 0)
			assert.False(t, greater.Equal(v))
		}
	}
}
