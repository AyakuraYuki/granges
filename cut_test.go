package granges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		}
	}
}
