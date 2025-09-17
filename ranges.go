package granges

import "fmt"

type Range[C Comparable] struct {
	lowerBound Cut[C]
	upperBound Cut[C]

	invalid bool
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

func create[C Comparable](lowerBound, upperBound Cut[C]) (Range[C], error) {
	if lowerBound.Compare(upperBound) > 0 ||
		lowerBound.cutType == AboveAll ||
		upperBound.cutType == BelowAll {

		return Invalid[C](), fmt.Errorf(
			"invalid range: %s..%s",
			lowerBound.DescribeAsLowerBound(),
			upperBound.DescribeAsUpperBound())

	}

	return Range[C]{lowerBound: lowerBound, upperBound: upperBound}, nil
}

// IsInvalid identifies the range is invalid or not
func (r Range[C]) IsInvalid() bool {
	return r.invalid
}

// HasLowerBound returns true if this range has a lower endpoint.
func (r Range[C]) HasLowerBound() bool {
	return r.lowerBound.cutType != BelowAll
}

// LowerEndpoint returns the lower endpoint of this range with ignoring
// ErrRangeSideUnbounded error.
func (r Range[C]) LowerEndpoint() C {
	c, _ := r.lowerBound.Endpoint()
	return c
}

// LowerEndpointE returns the lower endpoint of this range.
// If this range is unbounded below (that is, HasLowerBound returns false),
// the ErrRangeSideUnbounded will be returned.
func (r Range[C]) LowerEndpointE() (C, error) {
	return r.lowerBound.Endpoint()
}

// LowerBoundType returns the type of this range's lower bound:
// CLOSED if the range includes its lower endpoint, OPEN if it does not,
// Unbounded if this range is unbounded below (that is, HasLowerBound returns
// false).
func (r Range[C]) LowerBoundType() BoundType {
	t, _ := r.lowerBound.TypeAsLowerBound()
	return t
}

// LowerBoundTypeE returns the type of this range's lower bound:
// CLOSED if the range includes its lower endpoint, OPEN if it does not,
// Unbounded and ErrUnboundedCut error if this range is unbounded below (that
// is, HasLowerBound returns false).
func (r Range[C]) LowerBoundTypeE() (BoundType, error) {
	return r.lowerBound.TypeAsLowerBound()
}

// HasUpperBound returns true if this range has an upper endpoint.
func (r Range[C]) HasUpperBound() bool {
	return r.upperBound.cutType != AboveAll
}

// UpperEndpoint returns the upper endpoint of this range with ignoring
// ErrRangeSideUnbounded error.
func (r Range[C]) UpperEndpoint() C {
	c, _ := r.upperBound.Endpoint()
	return c
}

// UpperEndpointE returns the lower endpoint of this range.
// If this range is unbounded above (that is, HasUpperBound returns false),
// the ErrRangeSideUnbounded will be returned.
func (r Range[C]) UpperEndpointE() (C, error) {
	return r.upperBound.Endpoint()
}

// UpperBoundType returns the type of this range's upper bound:
// CLOSED if the range includes its upper endpoint, OPEN if it does not,
// Unbounded if this range is unbounded above (that is, HasUpperBound returns
// false).
func (r Range[C]) UpperBoundType() BoundType {
	t, _ := r.upperBound.TypeAsUpperBound()
	return t
}

// UpperBoundTypeE returns the type of this range's upper bound:
// CLOSED if the range includes its upper endpoint, OPEN if it does not,
// Unbounded and ErrUnboundedCut error if this range is unbounded above (that
// is, HasUpperBound returns false).
func (r Range[C]) UpperBoundTypeE() (BoundType, error) {
	return r.upperBound.TypeAsUpperBound()
}

// IsEmpty returns true if this range is of the form [v..v) or (v..v]. (This
// does not encompass ranges of the form (v..v), because such ranges are
// invalid and can't be constructed at all.)
func (r Range[C]) IsEmpty() bool {
	return r.lowerBound.Compare(r.upperBound) == 0
}

// Contains returns true if value is within the bounds of this range. For
// example, on the range [0..2), Contains(1) returns true, while Contains(2)
// returns false.
func (r Range[C]) Contains(value C) bool {
	return r.lowerBound.IsLessThan(value) && !r.upperBound.IsLessThan(value)
}

// ContainsAll returns true if every element in values is contained in this
// range.
func (r Range[C]) ContainsAll(values []C) bool {
	if len(values) == 0 {
		return true
	}
	for i := range values {
		if !r.Contains(values[i]) {
			return false
		}
	}
	return true
}

// Encloses returns true if the bounds of other do not extend outside the
// bounds of this range.
//
// Examples:
//   - [3..6] encloses [4..5]
//   - (3..6) encloses (3..6)
//   - [3..6] encloses [4..4) (even though the latter is empty)
//   - (3..6] does not enclose [3..6]
//   - [4..5] does not enclose (3..6) (even though it contains every value
//     contained by the latter range)
//   - [3..6] does not enclose (1..1] (even though it contains every value
//     contained by the latter range)
//
// Note that if a.Encloses(b), then b.Contains(v) implies a.Contains(v), but as
// the last two examples illustrate, the converse is not always true.
//
// Being reflexive, antisymmetric and transitive, the encloses relation defines
// a partial order over ranges. There exists a unique maximal range according
// to this relation, and also numerous minimal ranges. Enclosure also implies
// connectedness.
func (r Range[C]) Encloses(other Range[C]) bool {
	return r.lowerBound.Compare(other.lowerBound) <= 0 &&
		r.upperBound.Compare(other.upperBound) >= 0
}

// IsConnected returns true if there exists a (possibly empty) range which is
// enclosed by both this range and other.
//
// For example,
//
//   - [2, 4) and [5, 7) are not connected
//   - [2, 4) and [3, 5) are connected, because both enclose [3, 4)
//   - [2, 4) and [4, 6) are connected, because both enclose the empty range
//     [4, 4)
//
// Note that this range and other have a well-defined union and intersection
// (as a single, possibly-empty range) if and only if this method returns true.
//
// The connectedness relation is both reflexive and symmetric, but does not
// form an equivalence relation as it is not transitive.
func (r Range[C]) IsConnected(other Range[C]) bool {
	return r.lowerBound.Compare(other.upperBound) <= 0 &&
		other.lowerBound.Compare(r.upperBound) <= 0
}

// Intersection returns the maximal range enclosed by both this range and
// connectedRange, if such a range exists.
//
// An invalid range will be returned for disconnected ranges.
func (r Range[C]) Intersection(connectedRange Range[C]) Range[C] {
	intersection, _ := r.IntersectionE(connectedRange)
	return intersection
}

// IntersectionE returns the maximal range enclosed by both this range and
// connectedRange, if such a range exists.
//
// For example, the intersection of [1..5] and (3..7) is (3..5]. The resulting
// range may be empty; for example, [1..5) intersected with [5..7) yields the
// empty range [5..5).
//
// The intersection exists if and only if the two ranges are connected.
//
// The intersection operation is commutative, associative and idempotent, and
// its identity element is All.
//
// An error will be returned for disconnected ranges.
func (r Range[C]) IntersectionE(connectedRange Range[C]) (Range[C], error) {
	lowerCmp := r.lowerBound.Compare(connectedRange.lowerBound)
	upperCmp := r.upperBound.Compare(connectedRange.upperBound)
	if lowerCmp >= 0 && upperCmp <= 0 {
		return r, nil
	} else if lowerCmp <= 0 && upperCmp >= 0 {
		return connectedRange, nil
	} else {
		var (
			newLower Cut[C]
			newUpper Cut[C]
		)

		if lowerCmp >= 0 {
			newLower = r.lowerBound
		} else {
			newLower = connectedRange.lowerBound
		}

		if upperCmp <= 0 {
			newUpper = r.upperBound
		} else {
			newUpper = connectedRange.upperBound
		}

		if newLower.Compare(newUpper) > 0 {
			return Invalid[C](), fmt.Errorf(
				"intersection is undefined for disconnected ranges %s and %s",
				r, connectedRange)
		}

		return create(newLower, newUpper)
	}
}

// Gap returns the maximal range lying between this range and otherRange, if
// such a range exists. The resulting range may be empty if the two ranges are
// adjacent but non-overlapping.
//
// An invalid range will be returned if this range and otherRange have a
// nonempty intersection.
func (r Range[C]) Gap(other Range[C]) Range[C] {
	gap, _ := r.GapE(other)
	return gap
}

// GapE returns the maximal range lying between this range and otherRange, if
// such a range exists. The resulting range may be empty if the two ranges are
// adjacent but non-overlapping.
//
// For example, the gap of [1..5] and (7..10) is (5..7]. The resulting range
// may be empty; for example, the gap between [1..5) [5..7) yields the empty
// range [5..5).
//
// The gap exists if and only if the two ranges are either disconnected or
// immediately adjacent (any intersection must be an empty range).
//
// The gap operation is commutative.
//
// An error will be returned if this range and otherRange have a nonempty
// intersection.
func (r Range[C]) GapE(other Range[C]) (Range[C], error) {
	if r.lowerBound.Compare(other.upperBound) < 0 &&
		other.lowerBound.Compare(r.upperBound) < 0 {
		return Invalid[C](), fmt.Errorf(
			"ranges have a nonempty intersection: %s, %s", r, other)
	}

	var (
		isThisFirst = r.lowerBound.Compare(other.lowerBound) < 0
		first       Range[C]
		second      Range[C]
	)

	if isThisFirst {
		first, second = r, other
	} else {
		first, second = other, r
	}
	return create(first.upperBound, second.lowerBound)
}

// Span returns the minimal range that encloses both this range and other.
// For example, the span of [1..3] and (5..7) is [1..7).
//
// An invalid range will be returned if failed to create new range.
func (r Range[C]) Span(other Range[C]) Range[C] {
	span, _ := r.SpanE(other)
	return span
}

// SpanE returns the minimal range that encloses both this range and other.
// For example, the span of [1..3] and (5..7) is [1..7).
//
// If the input ranges are connected, the returned range can also be called
// their union. If they are not, note that the span might contain values that
// are not contained in either input range.
//
// Like intersection, this operation is commutative, associative and
// idempotent. Unlike it, it is always well-defined for any two input ranges.
//
// An error will be returned if failed to create new range.
func (r Range[C]) SpanE(other Range[C]) (Range[C], error) {
	lowerCmp := r.lowerBound.Compare(other.lowerBound)
	upperCmp := r.upperBound.Compare(other.upperBound)
	if lowerCmp <= 0 && upperCmp >= 0 {
		return r, nil
	} else if lowerCmp >= 0 && upperCmp <= 0 {
		return other, nil
	} else {
		var (
			newLower Cut[C]
			newUpper Cut[C]
		)
		if lowerCmp <= 0 {
			newLower = r.lowerBound
		} else {
			newLower = other.lowerBound
		}
		if upperCmp >= 0 {
			newUpper = r.upperBound
		} else {
			newUpper = other.upperBound
		}
		return create(newLower, newUpper)
	}
}

// Equal returns true if object is a range having the same endpoints and bound
// types as this range. Note that discrete ranges such as (1..4) and [2..3] are
// not equal to one another, despite the fact that they each contain precisely
// the same set of values. Similarly, empty ranges are not equal unless they
// have exactly the same representation, so [3..3), (3..3], (4..4] are all
// unequal.
func (r Range[C]) Equal(other Range[C]) bool {
	return r.lowerBound.Compare(other.lowerBound) == 0 &&
		r.upperBound.Compare(other.upperBound) == 0
}

func (r Range[C]) String() string {
	lowerStr := r.lowerBound.DescribeAsLowerBound()
	upperStr := r.upperBound.DescribeAsUpperBound()
	return fmt.Sprintf("%s..%s", lowerStr, upperStr)
}
