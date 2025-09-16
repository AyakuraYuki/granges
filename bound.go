package granges

// BoundType indicates whether an endpoint of some range is contained in the
// range itself ("closed") or not ("open"). If a range is unbounded on a side,
// it is neither open nor closed on that side; the bound simply does not exist.
type BoundType int

const Unbounded BoundType = -1

const (
	OPEN   BoundType = iota // open interval: ()
	CLOSED                  // closed interval: []
)
