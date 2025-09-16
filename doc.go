/*
Package granges implements a mathematical interval operation tool.

A range (or "interval") defines the boundaries around a contiguous span of
values of some Comparable type; for example, "integers from 1 to 100 inclusive".

# Types of ranges

Each end of the range may be bounded or unbounded. If bounded, there is an
associated endpoint value, and the range is considered to be either OPEN (does
not include the endpoint) or CLOSED (includes the endpoint) on that side. With
three possibilities on each side, this yields nine basic types of ranges,
enumerated below. (Notation: a square bracket ([]) indicates that the range is
CLOSED on that side; a parenthesis (()) means it is either open or unbounded.
The construct {x | statement} is read "the set of all x such that statement."

# Range Types

  - Open: (a..b) -> {x | a < x < b}
  - Closed: [a..b] -> {x | a <= x <= b}
  - OpenClosed: (a..b] -> {x | a < x <= b}
  - ClosedOpen: [a..b) -> {x | a <= x < b}
  - GreaterThan: (a..+∞) -> {x | x > a}
  - AtLeast: [a..+∞) -> {x | x >= a}
  - LessThan: (-∞..b) -> {x | x < b}
  - AtMost: (-∞..b] -> {x | x <= b}
  - All: (-∞..+∞) -> {x}

When both endpoints exist, the upper endpoint may not be less than the lower.
The endpoints may be equal only if at least one of the bounds is closed:

  - [a..a] : a singleton range
  - [a..a); (a..a] : empty ranges; also valid
  - (a..a) : invalid; an exception will be thrown

# Warnings

  - Use immutable value types only, if at all possible. If you must use a
    mutable type, do not allow the endpoint instances to mutate after the range
    is created!
  - Your value type's comparison method should be consistent with equals if at
    all possible. Otherwise, be aware that concepts used throughout this
    documentation such as "equal", "same", "unique" and so on actually refer
    to whether Compare returns zero, not whether equals returns true.

# Other notes

  - All ranges are shallow-immutable.
*/
package granges
