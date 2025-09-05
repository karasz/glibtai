// Copyright © 2018 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of glibtai, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>

// Package glibtai is a partial Go implementation of libtai. See
// http://cr.yp.to/libtai/ for more information.
package glibtai

import (
	"math"
	"testing"
	"time"
)

func TestConversion(t *testing.T) {
	z, err := TAIfromString("@40000000036db755")
	if err != nil {
		t.Error(err)
	}
	tt := TAITime(z)
	if !tt.Equal(time.Date(1971, time.October, 28, 18, 19, 55, 0, time.UTC)) {
		t.Error("Expected 1971-10-28 18:19:55 +0000 UTC, got ", tt)
	}
}

func TestTAIfromTime(t *testing.T) {
	tt := time.Date(2018, time.February, 14, 19, 31, 10, 0, time.UTC)
	z := TAIfromTime(tt)
	q := TAITime(z)
	if !tt.Equal(q) {
		t.Errorf("%v is not equal with %v", tt, q)
	}
}

func TestTAIAddOverflow(t *testing.T) {
	tests := []struct {
		name     string
		tai      TAI
		duration time.Duration
		expected uint64
	}{
		{
			name:     "near max value positive overflow",
			tai:      TAI{x: math.MaxUint64 - 5},
			duration: 10 * time.Second,
			expected: 4, // wraps around: (MaxUint64 - 5 + 10) % 2^64 = 4
		},
		{
			name:     "exact max value overflow",
			tai:      TAI{x: math.MaxUint64},
			duration: 1 * time.Second,
			expected: 0, // wraps to 0
		},
		{
			name:     "large positive duration",
			tai:      TAI{x: 1000},
			duration: time.Duration(math.MaxInt64),
			expected: 1000 + uint64(math.MaxInt64/int64(time.Second)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TAIAdd(tc.tai, tc.duration)
			if result.x != tc.expected {
				t.Errorf("Expected x=%d, got x=%d", tc.expected, result.x)
			}
		})
	}
}

func TestTAIAddUnderflow(t *testing.T) {
	tests := []struct {
		name     string
		tai      TAI
		duration time.Duration
		expected uint64
	}{
		{
			name:     "small value negative underflow",
			tai:      TAI{x: 5},
			duration: -10 * time.Second,
			expected: math.MaxUint64 - 4, // wraps around: 5 - 10 = -5 wraps to MaxUint64 - 4
		},
		{
			name:     "zero value underflow",
			tai:      TAI{x: 0},
			duration: -1 * time.Second,
			expected: math.MaxUint64, // 0 - 1 wraps to MaxUint64
		},
		{
			name:     "large negative duration",
			tai:      TAI{x: math.MaxUint64 / 2},
			duration: -time.Duration(math.MaxInt64),
			expected: math.MaxUint64/2 - uint64(math.MaxInt64/int64(time.Second)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TAIAdd(tc.tai, tc.duration)
			if result.x != tc.expected {
				t.Errorf("Expected x=%d, got x=%d", tc.expected, result.x)
			}
		})
	}
}

func TestTAIAddEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		tai      TAI
		duration time.Duration
		check    func(TAI, TAI) bool
	}{
		{
			name:     "zero duration",
			tai:      TAI{x: 12345},
			duration: 0,
			check: func(original, result TAI) bool {
				return original.x == result.x
			},
		},
		{
			name:     "sub-second precision loss",
			tai:      TAI{x: 1000},
			duration: 500 * time.Millisecond,
			check: func(original, result TAI) bool {
				return result.x == original.x // sub-second part should be truncated
			},
		},
		{
			name:     "math.MinInt64 edge case",
			tai:      TAI{x: math.MaxUint64 / 2},
			duration: time.Duration(math.MinInt64),
			check: func(original, result TAI) bool {
				// math.MinInt64 as duration = -9223372036854775808 nanoseconds
				// Converted to seconds: -9223372036 seconds (not MinInt64!)
				seconds := int64(time.Duration(math.MinInt64) / time.Second)
				expectedSubtract := uint64(-seconds) // -(-9223372036) = 9223372036
				expected := original.x - expectedSubtract
				return result.x == expected
			},
		},
		{
			name:     "math.MaxInt64 positive duration",
			tai:      TAI{x: 1000},
			duration: time.Duration(math.MaxInt64),
			check: func(original, result TAI) bool {
				// math.MaxInt64 nanoseconds = 9223372036 seconds
				expectedAdd := uint64(math.MaxInt64 / int64(time.Second))
				expected := original.x + expectedAdd
				return result.x == expected
			},
		},
		{
			name:     "extreme negative with wraparound",
			tai:      TAI{x: 100},
			duration: time.Duration(math.MinInt64),
			check: func(original, result TAI) bool {
				// Should wrap around due to uint64 underflow
				seconds := int64(time.Duration(math.MinInt64) / time.Second)
				expectedSubtract := uint64(-seconds)      // 9223372036
				expected := original.x - expectedSubtract // Will wrap around
				return result.x == expected
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TAIAdd(tc.tai, tc.duration)
			if !tc.check(tc.tai, result) {
				t.Errorf("Test failed for %s: original=%v, result=%v", tc.name, tc.tai, result)
			}
		})
	}
}
