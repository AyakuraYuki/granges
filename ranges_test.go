package granges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AyakuraYuki/granges"
)

func TestRange_IsConnected(t *testing.T) {
	tests := []struct {
		A, B granges.Range[int]
		Want bool
	}{
		{A: granges.Closed(3, 5), B: granges.Open(5, 6), Want: true},
		{A: granges.Closed(3, 5), B: granges.Closed(5, 6), Want: true},
		{A: granges.Closed(5, 6), B: granges.Closed(3, 5), Want: true},
		{A: granges.Closed(3, 5), B: granges.OpenClosed(5, 5), Want: true},
		{A: granges.Open(3, 5), B: granges.Closed(5, 6), Want: true},
		{A: granges.Closed(3, 7), B: granges.Open(6, 8), Want: true},
		{A: granges.Open(3, 7), B: granges.Closed(5, 6), Want: true},
		{A: granges.Closed(3, 5), B: granges.Closed(7, 8), Want: false},
		{A: granges.Closed(3, 5), B: granges.ClosedOpen(7, 7), Want: false},
	}

	for _, tt := range tests {
		get := tt.A.IsConnected(tt.B)
		if get != tt.Want {
			t.Errorf("IsConnected(%q, %q) = %v, want %v", tt.A, tt.B, get, tt.Want)
		}
	}
}

func TestRange_IsEmpty_1(t *testing.T) {
	r := granges.ClosedOpen(4, 4)
	assert.False(t, r.Contains(3))
	assert.False(t, r.Contains(4))
	assert.False(t, r.Contains(5))
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 4, r.LowerEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 4, r.UpperEndpoint())
	assert.EqualValues(t, granges.OPEN, r.UpperBoundType())
	assert.True(t, r.IsEmpty())
	assert.EqualValues(t, "[4..4)", r.String())
}

func TestRange_IsEmpty_2(t *testing.T) {
	r := granges.OpenClosed(4, 4)
	assert.False(t, r.Contains(3))
	assert.False(t, r.Contains(4))
	assert.False(t, r.Contains(5))
	assert.True(t, r.HasLowerBound())
	assert.EqualValues(t, 4, r.LowerEndpoint())
	assert.EqualValues(t, granges.OPEN, r.LowerBoundType())
	assert.True(t, r.HasUpperBound())
	assert.EqualValues(t, 4, r.UpperEndpoint())
	assert.EqualValues(t, granges.CLOSED, r.UpperBoundType())
	assert.True(t, r.IsEmpty())
	assert.EqualValues(t, "(4..4]", r.String())
}

func TestRange_ContainsAll(t *testing.T) {
	r := granges.Closed(3, 5)
	assert.True(t, r.ContainsAll([]int{3, 3, 4, 5}))
	assert.False(t, r.ContainsAll([]int{3, 3, 4, 5, 6}))

	var empty []int
	assert.True(t, granges.OpenClosed(3, 3).ContainsAll(empty))
}

func TestRange_Encloses_open(t *testing.T) {
	r := granges.Open(2, 5)

	assert.True(t, r.Encloses(r))
	assert.True(t, r.Encloses(granges.Open(2, 4)))
	assert.True(t, r.Encloses(granges.Open(3, 5)))
	assert.True(t, r.Encloses(granges.Open(3, 4)))

	assert.False(t, r.Encloses(granges.OpenClosed(2, 5)))
	assert.False(t, r.Encloses(granges.ClosedOpen(2, 5)))
	assert.False(t, r.Encloses(granges.Closed(1, 4)))
	assert.False(t, r.Encloses(granges.Closed(3, 6)))
	assert.False(t, r.Encloses(granges.GreaterThan(3)))
	assert.False(t, r.Encloses(granges.LessThan(3)))
	assert.False(t, r.Encloses(granges.AtLeast(3)))
	assert.False(t, r.Encloses(granges.AtMost(3)))
	assert.False(t, r.Encloses(granges.All[int]()))
}

func TestRange_Encloses_closed(t *testing.T) {
	r := granges.Closed(2, 5)

	assert.True(t, r.Encloses(r))
	assert.True(t, r.Encloses(granges.Open(2, 5)))
	assert.True(t, r.Encloses(granges.OpenClosed(2, 5)))
	assert.True(t, r.Encloses(granges.ClosedOpen(2, 5)))
	assert.True(t, r.Encloses(granges.Closed(3, 5)))
	assert.True(t, r.Encloses(granges.Closed(2, 4)))

	assert.False(t, r.Encloses(granges.Open(1, 6)))
	assert.False(t, r.Encloses(granges.GreaterThan(3)))
	assert.False(t, r.Encloses(granges.LessThan(3)))
	assert.False(t, r.Encloses(granges.AtLeast(3)))
	assert.False(t, r.Encloses(granges.AtMost(3)))
	assert.False(t, r.Encloses(granges.All[int]()))
}

func TestRange_Intersection_empty(t *testing.T) {
	r := granges.ClosedOpen(3, 3)
	ri, err := r.IntersectionE(r)
	assert.NoError(t, err)
	assert.True(t, r.Equal(ri))

	intersection, err := r.IntersectionE(granges.Open(3, 5))
	assert.ErrorContains(t, err, "connected")
	assert.True(t, intersection.IsInvalid())
	intersection, err = r.IntersectionE(granges.Closed(0, 2))
	assert.ErrorContains(t, err, "connected")
	assert.True(t, intersection.IsInvalid())
}

func TestRange_Intersection_deFactoEmpty(t *testing.T) {
	{
		r := granges.Open(3, 4)
		ri, err := r.IntersectionE(r)
		assert.NoError(t, err)
		assert.True(t, r.Equal(ri))

		assert.True(t, granges.OpenClosed(3, 3).Equal(r.Intersection(granges.AtMost(3))))
		assert.True(t, granges.ClosedOpen(4, 4).Equal(r.Intersection(granges.AtLeast(4))))

		intersection, err := r.IntersectionE(granges.LessThan(3))
		assert.ErrorContains(t, err, "connected")
		assert.True(t, intersection.IsInvalid())
		intersection, err = r.IntersectionE(granges.GreaterThan(4))
		assert.ErrorContains(t, err, "connected")
		assert.True(t, intersection.IsInvalid())
	}

	{
		r := granges.Closed(3, 4)
		ri, _ := r.IntersectionE(granges.GreaterThan(4))
		assert.True(t, granges.OpenClosed(4, 4).Equal(ri))
	}
}

func TestRange_Intersection_singleton(t *testing.T) {
	r := granges.Singleton(3)
	ri, err := r.IntersectionE(r)
	assert.NoError(t, err)
	assert.True(t, r.Equal(ri))

	assert.True(t, r.Equal(r.Intersection(granges.AtMost(4))))
	assert.True(t, r.Equal(r.Intersection(granges.AtMost(3))))
	assert.True(t, r.Equal(r.Intersection(granges.AtLeast(3))))
	assert.True(t, r.Equal(r.Intersection(granges.AtLeast(2))))

	assert.True(t, granges.ClosedOpen(3, 3).Equal(r.Intersection(granges.LessThan(3))))
	assert.True(t, granges.OpenClosed(3, 3).Equal(r.Intersection(granges.GreaterThan(3))))

	intersection, err := r.IntersectionE(granges.AtLeast(4))
	assert.ErrorContains(t, err, "connected")
	assert.True(t, intersection.IsInvalid())
	intersection, err = r.IntersectionE(granges.AtMost(2))
	assert.ErrorContains(t, err, "connected")
	assert.True(t, intersection.IsInvalid())
}

func TestRange_Intersection_general(t *testing.T) {
	r := granges.Closed(4, 8)

	// separate below
	_, err := r.IntersectionE(granges.Closed(0, 2))
	assert.ErrorContains(t, err, "connected")

	// adjacent below
	assert.True(t, granges.ClosedOpen(4, 4).Equal(r.Intersection(granges.ClosedOpen(2, 4))))

	// overlap below
	assert.True(t, granges.Closed(4, 6).Equal(r.Intersection(granges.Closed(2, 6))))

	// enclosed with same start
	assert.True(t, granges.Closed(4, 6).Equal(r.Intersection(granges.Closed(4, 6))))

	// enclosed, interior
	assert.True(t, granges.Closed(5, 7).Equal(r.Intersection(granges.Closed(5, 7))))

	// enclosed with same end
	assert.True(t, granges.Closed(6, 8).Equal(r.Intersection(granges.Closed(6, 8))))

	// equal
	assert.True(t, r.Equal(r.Intersection(r)))

	// enclosing with same start
	assert.True(t, r.Equal(r.Intersection(granges.Closed(4, 10))))

	// enclosing with same end
	assert.True(t, r.Equal(r.Intersection(granges.Closed(2, 8))))

	// enclosing, exterior
	assert.True(t, r.Equal(r.Intersection(granges.Closed(2, 10))))

	// overlap above
	assert.True(t, granges.Closed(6, 8).Equal(r.Intersection(granges.Closed(6, 10))))

	// adjacent above
	assert.True(t, granges.OpenClosed(8, 8).Equal(r.Intersection(granges.OpenClosed(8, 10))))

	// separate above
	intersection, err := r.IntersectionE(granges.Closed(10, 12))
	assert.ErrorContains(t, err, "connected")
	assert.True(t, intersection.IsInvalid())
}

func TestRange_Gap_overlapping(t *testing.T) {
	r := granges.ClosedOpen(3, 5)

	gap, err := r.GapE(granges.Closed(4, 6))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())

	gap, err = r.GapE(granges.Closed(2, 4))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())

	gap, err = r.GapE(granges.Closed(2, 3))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())
}

func TestRange_Gap_invalidRangesWithInfinity(t *testing.T) {
	gap, err := granges.AtLeast(1).GapE(granges.AtLeast(2))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())

	gap, err = granges.AtLeast(2).GapE(granges.AtLeast(1))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())

	gap, err = granges.AtMost(1).GapE(granges.AtMost(2))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())

	gap, err = granges.AtMost(2).GapE(granges.AtMost(1))
	assert.Error(t, err)
	assert.True(t, gap.IsInvalid())
}

func TestRange_Gap_connectedAdjacentYieldsEmpty(t *testing.T) {
	r := granges.Open(3, 4)
	assert.True(t, granges.ClosedOpen(4, 4).Equal(r.Gap(granges.AtLeast(4))))
	assert.True(t, granges.OpenClosed(3, 3).Equal(r.Gap(granges.AtMost(3))))
}

func TestRange_Gap_general(t *testing.T) {
	openRange := granges.Open(4, 8)
	closedRange := granges.Closed(4, 8)

	// first range open end, second range open start
	assert.True(t, granges.Closed(2, 4).Equal(granges.LessThan(2).Gap(openRange)))
	assert.True(t, granges.Closed(2, 4).Equal(openRange.Gap(granges.LessThan(2))))

	// first range closed end, second range open start
	assert.True(t, granges.OpenClosed(2, 4).Equal(granges.AtMost(2).Gap(openRange)))
	assert.True(t, granges.OpenClosed(2, 4).Equal(openRange.Gap(granges.AtMost(2))))

	// first range open end, second range closed start
	assert.True(t, granges.ClosedOpen(2, 4).Equal(granges.LessThan(2).Gap(closedRange)))
	assert.True(t, granges.ClosedOpen(2, 4).Equal(closedRange.Gap(granges.LessThan(2))))

	// first range closed end, second range closed start
	assert.True(t, granges.Open(2, 4).Equal(granges.AtMost(2).Gap(closedRange)))
	assert.True(t, granges.Open(2, 4).Equal(closedRange.Gap(granges.AtMost(2))))
}

func TestRange_Span(t *testing.T) {
	r := granges.Closed(4, 8)

	// separate below
	assert.True(t, granges.Closed(0, 8).Equal(r.Span(granges.Closed(0, 2))))
	assert.True(t, granges.AtMost(8).Equal(r.Span(granges.AtMost(2))))

	// adjacent below
	assert.True(t, granges.Closed(2, 8).Equal(r.Span(granges.ClosedOpen(2, 4))))
	assert.True(t, granges.AtMost(8).Equal(r.Span(granges.LessThan(4))))

	// overlap below
	assert.True(t, granges.Closed(2, 8).Equal(r.Span(granges.Closed(2, 6))))
	assert.True(t, granges.AtMost(8).Equal(r.Span(granges.AtMost(6))))

	// enclosed with same start
	assert.True(t, r.Equal(r.Span(granges.Closed(4, 6))))

	// enclosed, interior
	assert.True(t, r.Equal(r.Span(granges.Closed(5, 7))))

	// enclosed with same end
	assert.True(t, r.Equal(r.Span(granges.Closed(6, 8))))

	// equal
	assert.True(t, r.Equal(r.Span(r)))

	// enclosing with same start
	assert.True(t, granges.Closed(4, 10).Equal(r.Span(granges.Closed(4, 10))))
	assert.True(t, granges.AtLeast(4).Equal(r.Span(granges.AtLeast(4))))

	// enclosing with same end
	assert.True(t, granges.Closed(2, 8).Equal(r.Span(granges.Closed(2, 8))))
	assert.True(t, granges.AtMost(8).Equal(r.Span(granges.AtMost(8))))

	// enclosing, exterior
	assert.True(t, granges.Closed(2, 10).Equal(r.Span(granges.Closed(2, 10))))
	assert.True(t, granges.All[int]().Equal(r.Span(granges.All[int]())))

	// overlap above
	assert.True(t, granges.Closed(4, 10).Equal(r.Span(granges.Closed(6, 10))))
	assert.True(t, granges.AtLeast(4).Equal(r.Span(granges.AtLeast(6))))

	// adjacent above
	assert.True(t, granges.Closed(4, 10).Equal(r.Span(granges.OpenClosed(8, 10))))
	assert.True(t, granges.AtLeast(4).Equal(r.Span(granges.GreaterThan(8))))

	// separate above
	assert.True(t, granges.Closed(4, 12).Equal(r.Span(granges.Closed(10, 12))))
	assert.True(t, granges.AtLeast(4).Equal(r.Span(granges.AtLeast(10))))
}

func TestRange_Equal(t *testing.T) {
	assert.True(t, granges.Open(1, 5).Equal(granges.New(1, granges.OPEN, 5, granges.OPEN)))
	assert.True(t, granges.GreaterThan(2).Equal(granges.GreaterThan(2)))
	assert.True(t, granges.All[int]().Equal(granges.All[int]()))

	upTo, _ := granges.UpTo(7, granges.CLOSED)
	assert.True(t, granges.AtMost(7).Equal(upTo))

	upTo, _ = granges.UpTo(7, granges.OPEN)
	assert.True(t, granges.LessThan(7).Equal(upTo))

	// with invalid bound type
	upTo, _ = granges.UpTo(7, granges.Unbounded)
	assert.True(t, upTo.IsInvalid())

	downTo, _ := granges.DownTo(1, granges.CLOSED)
	assert.True(t, granges.AtLeast(1).Equal(downTo))

	downTo, _ = granges.DownTo(1, granges.OPEN)
	assert.True(t, granges.GreaterThan(1).Equal(downTo))

	// with invalid bound type
	downTo, _ = granges.DownTo(1, granges.Unbounded)
	assert.True(t, downTo.IsInvalid())

	assert.True(t, granges.Open(1, 7).Equal(granges.New(1, granges.OPEN, 7, granges.OPEN)))
	assert.True(t, granges.OpenClosed(1, 7).Equal(granges.New(1, granges.OPEN, 7, granges.CLOSED)))
	assert.True(t, granges.Closed(1, 7).Equal(granges.New(1, granges.CLOSED, 7, granges.CLOSED)))
	assert.True(t, granges.ClosedOpen(1, 7).Equal(granges.New(1, granges.CLOSED, 7, granges.OPEN)))
}
