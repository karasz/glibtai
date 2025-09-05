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

func TestTAINfromTime(t *testing.T) {
	tt := time.Date(2018, time.February, 14, 19, 31, 10, 0, time.UTC)
	z := TAINfromTime(tt)
	q := TAINTime(z)
	if !tt.Equal(q) {
		t.Errorf("%v is not equal with %v", tt, q)
	}
}

func TestTAINAddOverflow(t *testing.T) {
	tests := []struct {
		name     string
		tain     TAIN
		duration time.Duration
		expected TAIN
	}{
		{
			name:     "seconds overflow wraps around",
			tain:     TAIN{sec: math.MaxUint64 - 2, nano: 500000000},
			duration: 5 * time.Second,
			expected: TAIN{sec: 2, nano: 500000000}, // (MaxUint64 - 2 + 5) % 2^64 = 2
		},
		{
			name:     "nanoseconds carry causes seconds overflow",
			tain:     TAIN{sec: math.MaxUint64, nano: 500000000},
			duration: 600 * time.Millisecond,
			expected: TAIN{sec: 0, nano: 100000000}, // sec wraps to 0, nano = 500M + 600M - 1000M = 100M
		},
		{
			name:     "large duration overflow",
			tain:     TAIN{sec: math.MaxUint64 / 2, nano: 0},
			duration: time.Duration(math.MaxInt64),
			expected: TAIN{
				sec:  math.MaxUint64/2 + uint64(math.MaxInt64/int64(time.Second)),
				nano: uint32(math.MaxInt64 % int64(time.Second)),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TAINAdd(tc.tain, tc.duration)
			if result.sec != tc.expected.sec || result.nano != tc.expected.nano {
				t.Errorf("Expected sec=%d nano=%d, got sec=%d nano=%d",
					tc.expected.sec, tc.expected.nano, result.sec, result.nano)
			}
		})
	}
}

func checkTAINExpected(t *testing.T, expected TAIN, result TAIN) {
	if result.sec != expected.sec || result.nano != expected.nano {
		t.Errorf("Expected sec=%d nano=%d, got sec=%d nano=%d",
			expected.sec, expected.nano, result.sec, result.nano)
	}
}

func checkTAINWithFunc(t *testing.T, name string, original TAIN, result TAIN, checkFunc func(TAIN, TAIN) bool) {
	if !checkFunc(original, result) {
		t.Errorf("Check failed for %s: original=sec:%d nano:%d, result=sec:%d nano:%d",
			name, original.sec, original.nano, result.sec, result.nano)
	}
}

func checkLargeNegativeDuration(original, result TAIN) bool {
	totalNanos := int64(-time.Duration(math.MaxInt64))
	addSecs := totalNanos / 1e9
	addNanos := totalNanos % 1e9

	expectedSec := original.sec - uint64(-addSecs)
	expectedNano := int64(original.nano) + addNanos

	if expectedNano < 0 {
		expectedSec--
		expectedNano += 1e9
	}

	return result.sec == expectedSec && result.nano == uint32(expectedNano)
}

type tainTestCase struct {
	name     string
	tain     TAIN
	duration time.Duration
	expected *TAIN
	check    func(TAIN, TAIN) bool
}

func runTAINTests(t *testing.T, tests []tainTestCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TAINAdd(tc.tain, tc.duration)
			runSingleTAINTest(t, tc, result)
		})
	}
}

func runSingleTAINTest(t *testing.T, tc tainTestCase, result TAIN) {
	if tc.expected != nil {
		checkTAINExpected(t, *tc.expected, result)
	} else if tc.check != nil {
		checkTAINWithFunc(t, tc.name, tc.tain, result, tc.check)
	}
}

func TestTAINAddUnderflow(t *testing.T) {
	tests := []tainTestCase{
		{
			name:     "seconds underflow wraps around",
			tain:     TAIN{sec: 2, nano: 500000000},
			duration: -5 * time.Second,
			expected: &TAIN{sec: math.MaxUint64 - 2, nano: 500000000},
		},
		{
			name:     "nanoseconds borrow causes seconds underflow",
			tain:     TAIN{sec: 0, nano: 300000000},
			duration: -500 * time.Millisecond,
			expected: &TAIN{sec: math.MaxUint64, nano: 800000000},
		},
		{
			name:     "zero seconds underflow",
			tain:     TAIN{sec: 0, nano: 0},
			duration: -1 * time.Second,
			expected: &TAIN{sec: math.MaxUint64, nano: 0},
		},
		{
			name:     "large negative duration with borrow",
			tain:     TAIN{sec: math.MaxUint64 / 2, nano: 999999999},
			duration: -time.Duration(math.MaxInt64),
			check:    checkLargeNegativeDuration,
		},
	}

	runTAINTests(t, tests)
}

func checkZeroDuration(original, result TAIN) bool {
	return original.sec == result.sec && original.nano == result.nano
}

func checkNanosecondBoundary(original, result TAIN) bool {
	return result.sec == original.sec+1 && result.nano == 0
}

func checkNegativeNanosecondBoundary(original, result TAIN) bool {
	return result.sec == original.sec-1 && result.nano == 999999999
}

func checkSubNanosecondPrecision(original, result TAIN) bool {
	expectedNano := original.nano + 123456789
	expectedSec := original.sec
	if expectedNano >= 1000000000 {
		expectedSec++
		expectedNano -= 1000000000
	}
	return result.sec == expectedSec && result.nano == expectedNano
}

func TestTAINAddEdgeCases(t *testing.T) {
	tests := []tainTestCase{
		{
			name:     "zero duration no change",
			tain:     TAIN{sec: 12345, nano: 678900000},
			duration: 0,
			check:    checkZeroDuration,
		},
		{
			name:     "exact nanosecond boundary",
			tain:     TAIN{sec: 1000, nano: 999999999},
			duration: 1 * time.Nanosecond,
			check:    checkNanosecondBoundary,
		},
		{
			name:     "negative nanosecond boundary",
			tain:     TAIN{sec: 1000, nano: 0},
			duration: -1 * time.Nanosecond,
			check:    checkNegativeNanosecondBoundary,
		},
		{
			name:     "sub-nanosecond precision preserved",
			tain:     TAIN{sec: 1000, nano: 500000000},
			duration: 123456789 * time.Nanosecond,
			check:    checkSubNanosecondPrecision,
		},
	}

	runTAINTests(t, tests)
}
