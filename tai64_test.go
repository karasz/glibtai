// Copyright © 2018 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of glibtai, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>

// Package glibtai is a partial Go implementation of libtai. See
// http://cr.yp.to/libtai/ for more information.
package glibtai

import (
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
