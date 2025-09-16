# glibtai

[![Go Report Card](https://goreportcard.com/badge/github.com/karasz/glibtai)](
https://goreportcard.com/report/github.com/karasz/glibtai)
[![Unlicensed](https://img.shields.io/badge/license-Unlicense-blue.svg)](
https://github.com/karasz/gnocco/blob/master/UNLICENSE)
[![Status](https://godoc.org/github.com/karasz/glibtai?status.svg)](
http://godoc.org/github.com/karasz/glibtai)

A pure Go implementation of TAI64, TAI64N, and TAI64NA timestamps as
specified by D. J. Bernstein's [libtai](https://cr.yp.to/libtai.html).

## Features

- **TAI64**: 64-bit TAI timestamps with 1-second precision
- **TAI64N**: 96-bit TAI timestamps with nanosecond precision
- **Leap second handling**: Automatic UTC to TAI conversion with built-in
  leap second table
- **High performance**: Optimized arithmetic operations with comprehensive
  benchmarks
- **Overflow safety**: Proper wraparound semantics following libtai's
  modulo 2^64 design
- **String formatting**: Standard `@` prefixed hexadecimal representation
- **Go integration**: Native `time.Time` conversion and `time.Duration`
  arithmetic

## Installation

```bash
go get github.com/karasz/glibtai
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/karasz/glibtai"
)

func main() {
    // Get current TAI64 timestamp
    tai := glibtai.TAINow()
    fmt.Printf("Current TAI64: %s\n", tai)

    // Add duration
    future := glibtai.TAIAdd(tai, 1*time.Hour)
    fmt.Printf("One hour later: %s\n", future)

    // Convert to Go time
    goTime := glibtai.TAITime(tai)
    fmt.Printf("As Go time: %s\n", goTime)

    // Work with nanosecond precision
    tain := glibtai.TAINNow()
    fmt.Printf("TAI64N: %s\n", tain)
}
```

## API Reference

### TAI64 Functions

#### Core Operations

```go
// Get current timestamp
tai := TAINow()                           // Current TAI64 timestamp

// Time conversion
tai := TAIfromTime(time.Now())           // From time.Time
goTime := TAITime(tai)                   // To time.Time

// Arithmetic
future := TAIAdd(tai, time.Hour)         // Add duration
diff, err := TAISub(tai2, tai1)         // Subtract timestamps
```

#### String Operations

```go
// String conversion
str := tai.String()                      // "@40000000036DB755"
tai, err := TAIfromString(str)           // Parse from string

// Binary serialization
bytes := TAIPack(tai)                    // 8-byte big-endian
tai := TAIUnpack(bytes)                  // Unpack from bytes
```

### TAI64N Functions

#### TAI64N Core Operations

```go
// Get current nanosecond-precision timestamp
tain := TAINNow()                        // Current TAI64N timestamp

// Time conversion
tain := TAINfromTime(time.Now())         // From time.Time
goTime := TAINTime(tain)                 // To time.Time

// Arithmetic (nanosecond-aware)
future := TAINAdd(tain, time.Nanosecond) // Add duration with carry
diff, err := TAINSub(tain2, tain1)      // Subtract with borrow
```

#### TAI64N String Operations

```go
// String conversion (24 hex chars)
str := tain.String()                     // "@40000000036DB755AB4CDE12"
tain, err := TAINfromString(str)         // Parse from string

// Binary serialization
bytes := TAINPack(tain)                  // 12-byte big-endian
tain := TAINUnpack(bytes)                // Unpack from bytes
```

## Time Scales and Precision

| Format | Size | Precision | Range | Use Case |
|--------|------|-----------|-------|----------|
| TAI64  | 8 bytes | 1 second | ~584 billion years | Log timestamps, general use |
| TAI64N | 12 bytes | 1 nanosecond | ~584 billion years | High-precision timing |

## Leap Seconds

This library handles leap seconds automatically using a built-in table updated
through 2017. The conversion accounts for the difference between UTC
(with leap seconds) and TAI (monotonic atomic time).

**Current UTC-TAI offset:** 37 seconds (as of 2017)

**Note:** For timestamps after 2017, you may need to update the leap second
table if new leap seconds are announced.

## Performance

Run benchmarks to see performance characteristics:

```bash
make bench
```

Typical performance on modern hardware:

- TAI64 operations: ~1-2 ns/op
- TAI64N operations: ~2-4 ns/op
- String operations: ~100-400 ns/op (before optimization)

## Error Handling

The library follows Go conventions:

```go
// Operations that can fail return error
tai, err := TAIfromString("@invalid")
if err != nil {
    log.Fatal(err)
}

// Arithmetic operations use overflow-safe wraparound
tai := TAI{x: math.MaxUint64}
overflowed := TAIAdd(tai, time.Second)  // Wraps to 0, no error
```

## Examples

### Basic Usage

```go
// Create timestamp from current time
tai := TAIfromTime(time.Now())

// Format as string
fmt.Printf("TAI64: %s\n", tai.String())

// Add 30 minutes
later := TAIAdd(tai, 30*time.Minute)
```

### High-Precision Timing

```go
// Nanosecond precision timing
start := TAINNow()
// ... do work ...
end := TAINNow()

duration, err := TAINSub(end, start)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Operation took: %v\n", duration)
```

### Working with Historical Dates

```go
// Create timestamp for Unix epoch
epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
tai := TAIfromTime(epoch)

// Account for leap seconds automatically
fmt.Printf("Unix epoch in TAI: %s\n", tai.String())
```

### Binary Protocol Integration

```go
// Pack timestamp for network transmission
tai := TAINow()
data := TAIPack(tai)

// Send data over network...

// Unpack on receiver
received := TAIUnpack(data)
fmt.Printf("Received timestamp: %s\n", received.String())
```

## Testing

The library includes comprehensive tests covering:

- Overflow and underflow scenarios
- Edge cases with extreme values
- Leap second calculations
- String parsing and formatting
- Cross-platform compatibility

```bash
go test -v                    # Run all tests
go test -bench=.             # Run benchmarks
make test                    # Full test suite
```

## Compatibility

- **Go version:** 1.16+
- **Architecture:** All platforms supported by Go
- **libtai compatibility:** Follows DJB's specification exactly

## References

- [TAI64 specification](https://cr.yp.to/libtai/tai64.html) by
  D. J. Bernstein
- [libtai library](https://cr.yp.to/libtai.html) - Original C implementation
- [International Atomic Time (TAI)][tai-wiki]

[tai-wiki]: https://en.wikipedia.org/wiki/International_Atomic_Time

## TODO

The following API functions from the original DJB libtai implementation are not yet implemented in this Go port:

### High-Precision Time (TAIA - TAI64NA)
- **TAIA support**: TAI64NA (1-attosecond precision) time format and operations
  - `taia_now()` - Get current time with attosecond precision
  - `taia_add()`, `taia_sub()` - Arithmetic operations for TAIA
  - `taia_pack()`, `taia_unpack()` - Serialization for 16-byte TAIA format
  - `taia_less()` - Comparison operations
  - `taia_half()` - Divide time by 2
  - `taia_approx()` - Convert to floating-point approximation
  - `taia_frac()` - Extract fractional part
  - `taia_fmtfrac()` - Format fractional seconds

### Calendar Date Operations (caldate)
- **Calendar date handling**: Year-month-day operations and conversions
  - `caldate_frommjd()` - Convert from Modified Julian Day number
  - `caldate_mjd()` - Convert to Modified Julian Day number
  - `caldate_normalize()` - Normalize invalid dates (e.g., Feb 30 → Mar 2)
  - `caldate_fmt()`, `caldate_scan()` - String formatting and parsing
  - `caldate_easter()` - Calculate Easter date for given year

### Calendar Time Operations (caltime)
- **Calendar time with timezone**: Complete date/time with UTC offset
  - `caltime_tai()` - Convert from TAI to calendar time in UTC
  - `caltime_utc()` - Convert from calendar time to TAI
  - `caltime_fmt()`, `caltime_scan()` - String formatting and parsing

### Leap Second Management (leapsecs)
- **Advanced leap second handling**: Beyond the current basic implementation
  - `leapsecs_init()` - Initialize leap second table
  - `leapsecs_read()` - Read leap second data from file
  - `leapsecs_add()` - Add leap seconds to TAI time
  - `leapsecs_sub()` - Subtract leap seconds from TAI time

### Missing Comparison and Utility Functions
- **TAI64/TAI64N comparisons**: `tai_less()`, `tain_less()` for time ordering
- **Validation functions**: Input validation for packed formats
- **Extended arithmetic**: More comprehensive overflow handling

**Priority**: Calendar operations (caldate/caltime) would provide the most value for general-purpose time handling, followed by TAIA support for ultra-high-precision applications.

## License

This software is released into the public domain. See [UNLICENSE](UNLICENSE)
for details.
