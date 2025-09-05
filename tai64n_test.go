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

func TestTAINfromTime(t *testing.T) {
	tt := time.Date(2018, time.February, 14, 19, 31, 10, 0, time.UTC)
	z := TAINfromTime(tt)
	q := TAINTime(z)
	if !tt.Equal(q) {
		t.Errorf("%v is not equal with %v", tt, q)
	}
}
