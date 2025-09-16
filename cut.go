package granges

import "fmt"

type CutType int

const (
	BelowAll   CutType = iota // -∞
	AboveAll                  // +∞
	BelowValue                // <= value
	AboveValue                // > value
)

// Cut is the implementation detail for the internal structure of Range
// instances. Represents a unique way of "cutting" a "number line" (actually
// of instances of type C, not necessarily "numbers") into two sections; this
// can be done below a certain value, above a certain value, below all values
// or above all values. With this object defined in this way, an interval can
// always be represented by a pair of Cut instances.
type Cut[C Comparable] struct {
	cutType  CutType
	endpoint C
}

func NewBelowAll[C Comparable]() Cut[C] {
	return Cut[C]{cutType: BelowAll}
}

func NewAboveAll[C Comparable]() Cut[C] {
	return Cut[C]{cutType: AboveAll}
}

func NewBelowValue[C Comparable](value C) Cut[C] {
	return Cut[C]{cutType: BelowValue, endpoint: value}
}

func NewAboveValue[C Comparable](value C) Cut[C] {
	return Cut[C]{cutType: AboveValue, endpoint: value}
}

func (c Cut[C]) Endpoint() (endpoint C, err error) {
	if c.cutType == BelowAll || c.cutType == AboveAll {
		return endpoint, ErrRangeSideUnbounded
	}
	return c.endpoint, nil
}

func (c Cut[C]) IsLessThan(value C) bool {
	switch c.cutType {
	case BelowAll:
		return true
	case AboveAll:
		return false
	case BelowValue:
		// if endpoint <= value, cut is on the left side or same of value
		return c.endpoint <= value
	case AboveValue:
		// if endpoint < value, cut is on the left side of value
		return c.endpoint < value
	default:
		return false
	}
}

func (c Cut[C]) TypeAsLowerBound() (BoundType, error) {
	switch c.cutType {
	case BelowValue:
		return CLOSED, nil
	case AboveValue:
		return OPEN, nil
	default:
		return Unbounded, ErrUnboundedCut
	}
}

func (c Cut[C]) TypeAsUpperBound() (BoundType, error) {
	switch c.cutType {
	case BelowValue:
		return OPEN, nil
	case AboveValue:
		return CLOSED, nil
	default:
		return Unbounded, ErrUnboundedCut
	}
}

func (c Cut[C]) DescribeAsLowerBound() string {
	switch c.cutType {
	case BelowAll:
		return "(-∞"
	case BelowValue:
		return fmt.Sprintf("[%v", c.endpoint)
	case AboveValue:
		return fmt.Sprintf("(%v", c.endpoint)
	default:
		return ""
	}
}

func (c Cut[C]) DescribeAsUpperBound() string {
	switch c.cutType {
	case AboveAll:
		return "+∞)"
	case BelowValue:
		return fmt.Sprintf("%v)", c.endpoint)
	case AboveValue:
		return fmt.Sprintf("%v]", c.endpoint)
	default:
		return ""
	}
}

func (c Cut[C]) Compare(other Cut[C]) int {
	// INF
	if c.cutType == BelowAll {
		if other.cutType == BelowAll {
			return 0 // same INF
		}
		return -1
	}
	if c.cutType == AboveAll {
		if other.cutType == AboveAll {
			return 0 // same INF
		}
		return 1
	}
	if other.cutType == BelowAll {
		return 1
	}
	if other.cutType == AboveAll {
		return -1
	}

	// compare endpoint
	if c.endpoint < other.endpoint {
		return -1
	}
	if c.endpoint > other.endpoint {
		return 1
	}

	// compare cut type in same endpoint
	if c.cutType == BelowValue && other.cutType == AboveValue {
		return -1
	}
	if c.cutType == AboveValue && other.cutType == BelowValue {
		return 1
	}

	return 0
}

func (c Cut[C]) Equal(other Cut[C]) bool {
	return c.Compare(other) == 0
}
