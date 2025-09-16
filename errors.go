package granges

import "errors"

var (
	ErrRangeSideUnbounded = errors.New("range unbounded on this side")
	ErrUnboundedCut       = errors.New("unbounded cut")
	ErrWrongBoundType     = errors.New("unknown bound type")
)
