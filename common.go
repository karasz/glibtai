// Copyright © 2018 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of glibtai, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>

// Package glibtai is a partial Go implementation of libtai. See
// http://cr.yp.to/libtai/ for more information.
package glibtai

// TAICONST is 2^62+10 representing the TAI label of the second Unix started
// 1970-01-01 00:00:00 +0000 UTC
const TAICONST = uint64(4611686018427387914)

// TAILength is the length of a TAI timestamp in bytes
const TAILength = 8

// TAINLength is the length of a TAIN timestamp in bytes
const TAINLength = 12
