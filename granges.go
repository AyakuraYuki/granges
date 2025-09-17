package granges

// Open returns a range that contains all values strictly greater than lower
// and strictly less than upper.
//
//	(lower..upper) = {x | lower < x < upper}
//
// An invalid range will be returned if lower is greater than or equal to upper.
func Open[C Comparable](lower, upper C) Range[C] {
	r, _ := OpenE(lower, upper)
	return r
}

// OpenE returns a range that contains all values strictly greater than lower
// and strictly less than upper.
//
//	(lower..upper) = {x | lower < x < upper}
//
// An invalid range with an error will be returned if lower is greater than or
// equal to upper.
func OpenE[C Comparable](lower, upper C) (Range[C], error) {
	return create(NewAboveValue(lower), NewBelowValue(upper))
}

// Closed returns a range that contains all values greater than or equal to
// lower and less than or equal to upper.
//
//	[lower..upper] = {x | lower <= x <= upper}
//
// An invalid range will be returned if lower is greater than upper.
func Closed[C Comparable](lower, upper C) Range[C] {
	r, _ := ClosedE(lower, upper)
	return r
}

// ClosedE returns a range that contains all values greater than or equal to
// lower and less than or equal to upper.
//
//	[lower..upper] = {x | lower <= x <= upper}
//
// An invalid range with an error will be returned if lower is greater than
// upper.
func ClosedE[C Comparable](lower, upper C) (Range[C], error) {
	return create(NewBelowValue(lower), NewAboveValue(upper))
}

// ClosedOpen returns a range that contains all values greater than or equal
// to lower and strictly less than upper.
//
//	[lower..upper) = {x | lower <= x < upper}
//
// An invalid range will be returned if lower is greater than upper.
func ClosedOpen[C Comparable](lower, upper C) Range[C] {
	r, _ := ClosedOpenE(lower, upper)
	return r
}

// ClosedOpenE returns a range that contains all values greater than or equal
// to lower and strictly less than upper.
//
//	[lower..upper) = {x | lower <= x < upper}
//
// An invalid range with an error will be returned if lower is greater than
// upper.
func ClosedOpenE[C Comparable](lower, upper C) (Range[C], error) {
	return create(NewBelowValue(lower), NewBelowValue(upper))
}

// OpenClosed returns a range that contains all values strictly greater than
// lower and less than or equal to upper.
//
//	(lower..upper] = {x | lower < x <= upper}
//
// An invalid range will be returned if lower is greater than upper.
func OpenClosed[C Comparable](lower, upper C) Range[C] {
	r, _ := OpenClosedE(lower, upper)
	return r
}

// OpenClosedE returns a range that contains all values strictly greater than
// lower and less than or equal to upper.
//
//	(lower..upper] = {x | lower < x <= upper}
//
// An invalid range with an error will be returned if lower is greater than
// upper.
func OpenClosedE[C Comparable](lower, upper C) (Range[C], error) {
	return create(NewAboveValue(lower), NewAboveValue(upper))
}

// New returns a range that contains any value from lower to upper, where each
// endpoint may be either inclusive (closed) or exclusive (open).
//
// An invalid range will be returned if lower is greater than upper.
func New[C Comparable](lower C, lowerType BoundType, upper C, upperType BoundType) Range[C] {
	r, _ := NewE(lower, lowerType, upper, upperType)
	return r
}

// NewE returns a range that contains any value from lower to upper, where each
// endpoint may be either inclusive (closed) or exclusive (open).
//
// An invalid range with an error will be returned if lower is greater than
// upper.
func NewE[C Comparable](lower C, lowerType BoundType, upper C, upperType BoundType) (Range[C], error) {
	var (
		lowerBound Cut[C]
		upperBound Cut[C]
	)
	if lowerType == OPEN {
		lowerBound = NewAboveValue(lower)
	} else {
		lowerBound = NewBelowValue(lower)
	}
	if upperType == OPEN {
		upperBound = NewBelowValue(upper)
	} else {
		upperBound = NewAboveValue(upper)
	}
	return create(lowerBound, upperBound)
}

// LessThan returns a range that contains all values strictly less than
// endpoint.
//
//	(-∞..upper) = {x | x < upper}
func LessThan[C Comparable](upper C) Range[C] {
	return Range[C]{lowerBound: NewBelowAll[C](), upperBound: NewBelowValue(upper)}
}

// AtMost returns a range that contains all values less than or equal to
// endpoint.
//
//	(-∞..upper] = {x | x <= upper}
func AtMost[C Comparable](upper C) Range[C] {
	return Range[C]{lowerBound: NewBelowAll[C](), upperBound: NewAboveValue(upper)}
}

// GreaterThan returns a range that contains all values strictly greater than
// endpoint.
//
//	(lower..+∞) = {x | lower < x}
func GreaterThan[C Comparable](lower C) Range[C] {
	return Range[C]{lowerBound: NewAboveValue(lower), upperBound: NewAboveAll[C]()}
}

// AtLeast returns a range that contains all values greater than or equal to
// endpoint.
//
//	[lower..+∞) = {x | lower <= x}
func AtLeast[C Comparable](lower C) Range[C] {
	return Range[C]{lowerBound: NewBelowValue(lower), upperBound: NewAboveAll[C]()}
}

// All returns a range that contains every value of type T.
//
//	(-∞..+∞) = {x}
func All[C Comparable]() Range[C] {
	return Range[C]{lowerBound: NewBelowAll[C](), upperBound: NewAboveAll[C]()}
}

// Singleton returns a Range that contains only the given value.
// The returned range is CLOSED on both ends.
//
//	(x) = {x}
func Singleton[C Comparable](value C) Range[C] {
	return Closed(value, value)
}

// UpTo returns a range with no lower bound up to the given endpoint, which
// may be either inclusive (closed) or exclusive (open).
// An empty range with an error will be return if wrong arguments received.
func UpTo[C Comparable](endpoint C, boundType BoundType) (Range[C], error) {
	switch boundType {
	case OPEN:
		return LessThan(endpoint), nil
	case CLOSED:
		return AtMost(endpoint), nil
	default:
		return Invalid[C](), ErrWrongBoundType
	}
}

// DownTo returns a range from the given endpoint, which may be either
// inclusive (closed) or exclusive (open), with no upper bound.
// An empty range with an error will be return if wrong arguments received.
func DownTo[C Comparable](endpoint C, boundType BoundType) (Range[C], error) {
	switch boundType {
	case OPEN:
		return GreaterThan(endpoint), nil
	case CLOSED:
		return AtLeast(endpoint), nil
	default:
		return Invalid[C](), ErrWrongBoundType
	}
}

// Invalid creates an explicitly invalid range of the specified comparable
// type.
//
// This method is useful when you need to represent the concept of
// "no valid range" or signal an error condition in contexts where returning a
// range is required but no meaningful range can be constructed.
func Invalid[C Comparable]() Range[C] {
	return Range[C]{invalid: true}
}
