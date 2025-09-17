# granges

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful mathematical interval operation library for Go, ported from Google Guava's Range.
`granges` provides comprehensive support for creating, manipulating, and performing operations on ranges (intervals) of comparable values.

## Features

- **Type-safe intervals** using Go generics for any comparable type
- **Nine types of ranges** including open, closed, half-open, and unbounded ranges
- **Rich set operations** including intersection, union (span), gap analysis, and containment checking
- **Immutable design** ensuring thread safety
- **Comprehensive API** with both error-returning and panic-free variants
- **Full test coverage** ensuring reliability

## Installation

```shell
go get -u github.com/AyakuraYuki/granges
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/AyakuraYuki/granges"
)

func main() {
	// Create a closed range [1, 10]
	r1 := granges.Closed(1, 10)
	fmt.Println(r1.Contains(5))
	fmt.Println(r1.Contains(15))
	// Output:
	// true
	// false

	// Create an open range (0, 20)
	r2 := granges.Open(0, 20)

	// Find intersection
	intersection := r1.Intersection(r2)
	fmt.Println(intersection)
	// Output:
	// [1, 10]

	// Create unbounded ranges
	atLeast5 := granges.AtLeast(5)
	lessThan100 := granges.LessThan(100)
	fmt.Println(atLeast5.String())
	fmt.Println(lessThan100.String())
	// Output:
	// [5, +∞)
	// (-∞, 100)
}

```

## Range Types

`granges` supports nine fundamental types of ranges:

| Type            | Notation   | Description        | Example                    |
|-----------------|------------|--------------------|----------------------------|
| **Open**        | `(a, b)`   | `{x \| a < x < b}` | `granges.Open(1, 5)`       |
| **Closed**      | `[a, b]`   | `{x \| a ≤ x ≤ b}` | `granges.Closed(1, 5)`     |
| **OpenClosed**  | `(a, b]`   | `{x \| a < x ≤ b}` | `granges.OpenClosed(1, 5)` |
| **ClosedOpen**  | `[a, b)`   | `{x \| a ≤ x < b}` | `granges.ClosedOpen(1, 5)` |
| **GreaterThan** | `(a, +∞)`  | `{x \| x > a}`     | `granges.GreaterThan(5)`   |
| **AtLeast**     | `[a, +∞)`  | `{x \| x ≥ a}`     | `granges.AtLeast(5)`       |
| **LessThan**    | `(-∞, b)`  | `{x \| x < b}`     | `granges.LessThan(10)`     |
| **AtMost**      | `(-∞, b]`  | `{x \| x ≤ b}`     | `granges.AtMost(10)`       |
| **All**         | `(-∞, +∞)` | `{x}`              | `granges.All[int]()`       |

### Supported Types

The library works with any comparable type:

- **Integers**: `int`, `int8`, `int16`, `int32`, `int64`
- **Unsigned integers**: `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- **Floating point**: `float32`, `float64`
- **Strings**: `string`

## Creating Ranges

### Basic Range Creation

```go
package main

import "github.com/AyakuraYuki/granges"

func main() {
	// Bounded ranges
	closed := granges.Closed(1, 10)        // [1, 10]
	open := granges.Open(1, 10)            // (1, 10)
	halfOpen1 := granges.ClosedOpen(1, 10) // [1, 10)
	halfOpen2 := granges.OpenClosed(1, 10) // (1, 10]

	// Unbounded ranges
	atLeast := granges.AtLeast(5)         // [5, +∞)
	greaterThan := granges.GreaterThan(5) // (5, +∞)
	atMost := granges.AtMost(10)          // (-∞, 10]
	lessThan := granges.LessThan(10)      // (-∞, 10)

	// Special ranges
	singleton := granges.Singleton(42) // [42, 42]
	all := granges.All[int]()          // (-∞, +∞)
}

```

### Advanced Range Creation

```go
package main

import (
	"log"

	"github.com/AyakuraYuki/granges"
)

func main() {
	// Custom bound types
	r := granges.New(1, granges.CLOSED, 10, granges.OPEN) // [1, 10)

	// Error-handling variants
	r, err := granges.ClosedE(10, 1) // Returns error for invalid range
	if err != nil {
		log.Fatal(err)
	}

	// Directional ranges
	upTo, _ := granges.UpTo(100, granges.CLOSED) // (-∞, 100]
	downTo, _ := granges.DownTo(0, granges.OPEN) // (0, +∞)
}

```

## Range Operations

### Containment Testing

```go
package main

import (
	"fmt"

	"github.com/AyakuraYuki/granges"
)

func main() {
	r := granges.Closed(1, 10)

	// Single value
	fmt.Println(r.Contains(5))  // true
	fmt.Println(r.Contains(15)) // false

	// Multiple values
	values := []int{2, 5, 8}
	fmt.Println(r.ContainsAll(values)) // true

	// Range enclosure
	inner := granges.Closed(3, 7)
	fmt.Println(r.Encloses(inner)) // true
}

```

### Set Operations

```go
package main

import (
	"fmt"

	"github.com/AyakuraYuki/granges"
)

func main() {
	r1 := granges.Closed(1, 10)
	r2 := granges.Closed(5, 15)

	// Intersection - overlapping part
	intersection := r1.Intersection(r2) // [5, 10]

	// Span - minimal range covering both
	span := r1.Span(r2) // [1, 15]

	// Gap - range between two non-overlapping ranges
	r3 := granges.Closed(20, 30)
	gap := r1.Gap(r3) // (10, 20)

	// Connectivity test
	fmt.Println(r1.IsConnected(r2)) // true
	fmt.Println(r1.IsConnected(r3)) // false
}
```

### Range Properties

```go
package main

import (
	"fmt"
	"log"

	"github.com/AyakuraYuki/granges"
)

func main() {
	r := granges.ClosedOpen(5, 10)

	// Bounds checking
	fmt.Println(r.HasLowerBound()) // true
	fmt.Println(r.HasUpperBound()) // true

	// Endpoint access
	fmt.Println(r.LowerEndpoint()) // 5
	fmt.Println(r.UpperEndpoint()) // 10

	// Bound types
	fmt.Println(r.LowerBoundType()) // CLOSED
	fmt.Println(r.UpperBoundType()) // OPEN

	// Special properties
	fmt.Println(r.IsEmpty()) // false

	// Safe endpoint access with error handling
	endpoint, err := r.LowerEndpointE()
	if err != nil {
		log.Printf("No lower bound: %v", err)
	}
}

```

## Error Handling

The library provides two API patterns:

1. **Simple API**: Returns zero values for errors, suitable when you're confident about input validity
2. **Error-aware API**: Methods ending with 'E' return errors for better error handling

```go
package main

import (
	"log"

	"github.com/AyakuraYuki/granges"
)

func main() {
	// Simple API - may return invalid ranges
	r1 := granges.Open(10, 5) // Invalid range, but no error returned
	if r1.IsInvalid() {
		log.Println("Invalid range detected")
	}

	// Error-aware API - explicit error handling
	r2, err := granges.OpenE(10, 5)
	if err != nil {
		log.Printf("Failed to create range: %v", err)
		return
	}
}

```

## Working with Different Types

```go
package main

import (
	"fmt"

	"github.com/AyakuraYuki/granges"
)

func main() {
	// Integer ranges
	intRange := granges.Closed(1, 100)

	// Float ranges
	floatRange := granges.Open(0.0, 1.0)
	fmt.Println(floatRange.Contains(0.5)) // true

	// String ranges (lexicographic ordering)
	stringRange := granges.Closed("apple", "orange")
	fmt.Println(stringRange.Contains("banana")) // true
	fmt.Println(stringRange.Contains("zebra"))  // false
}

```

## Advanced Examples

### Range Validation

```go
package main

import "github.com/AyakuraYuki/granges"

func validateScore(score int) bool {
	validRange := granges.Closed(0, 100)
	return validRange.Contains(score)
}

```

### Range Partitioning

```go
package main

import "github.com/AyakuraYuki/granges"

func categorizeAge(age int) string {
	child := granges.ClosedOpen(0, 13)
	teen := granges.ClosedOpen(13, 20)
	adult := granges.ClosedOpen(20, 65)
	senior := granges.AtLeast(65)

	switch {
	case child.Contains(age):
		return "child"
	case teen.Contains(age):
		return "teenager"
	case adult.Contains(age):
		return "adult"
	case senior.Contains(age):
		return "senior"
	default:
		return "invalid"
	}
}

```

### Range Merging

```go
package main

import "github.com/AyakuraYuki/granges"

func mergeOverlappingRanges(r1, r2 granges.Range[int]) (granges.Range[int], bool) {
	if r1.IsConnected(r2) {
		return r1.Span(r2), true
	}
	return granges.Invalid[int](), false // Return invalid range if not connected
}

```

## Thread Safety

All range operations are safe for concurrent use. Ranges are immutable after creation, so they can be safely shared between goroutines without additional synchronization.

```go
package main

import (
	"fmt"

	"github.com/AyakuraYuki/granges"
)

var sharedRange = granges.Closed(1, 100)

func worker(id int) {
	// Safe to use sharedRange concurrently
	if sharedRange.Contains(id) {
		fmt.Printf("Worker %d is in range\n", id)
	}
}

```

## Performance Considerations

- Range creation and operations are generally O(1)
- All ranges are lightweight value types
- No heap allocations for basic operations
- Comparable type constraints ensure efficient comparisons

## Best Practices

1. **Use appropriate range types**: Choose the most restrictive range type that fits your needs
2. **Validate inputs**: Use error-returning variants (`*E` methods) when dealing with user input
3. **Prefer immutable patterns**: Create new ranges instead of trying to modify existing ones
4. **Check connectivity**: Use `IsConnected()` before performing intersection operations
5. **Handle edge cases**: Be aware of empty ranges and unbounded ranges in your logic

## Error Types

The library defines several error types:

- `ErrRangeSideUnbounded`: Returned when trying to access endpoints of unbounded ranges
- `ErrUnboundedCut`: Returned when trying to get bound types of unbounded ranges
- `ErrWrongBoundType`: Returned when invalid bound types are provided

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This library is ported from Google Guava's Range implementation. Special thanks to the Guava team for the excellent design and comprehensive functionality that inspired this Go version.
