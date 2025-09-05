// Copyright © 2018 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of glibtai, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>

// Package glibtai is a partial Go implementation of libtai. See
// http://cr.yp.to/libtai/ for more information.
package glibtai

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// TAI struct to store a TAI timestamp
type TAI struct {
	x uint64
}

// TAINow returns the current timestamp in TAI struct
func TAINow() TAI {
	return TAI{x: TAICONST + lsoffset(time.Now()) + uint64(time.Now().Unix())}
}

// TAIAdd adds a time.Duration to a TAI timestamp
func TAIAdd(a TAI, b time.Duration) TAI {
	return TAI{x: a.x + uint64(b.Seconds())}
}

// TAISub subtracts two TAI timestamps
func TAISub(a, b TAI) (time.Duration, error) {
	x := a.x - b.x
	q, err := time.ParseDuration(fmt.Sprintf("%ds", x))
	return q, err
}

// TAITime returns a go time object from a TAI timestamp
func TAITime(t TAI) time.Time {
	tm := time.Unix(int64(t.x-TAICONST), 0).UTC()
	return tm.Add(-time.Duration(lsoffset(tm)) * time.Second)
}

// TAIPack packs a TAI timestamp into a byte array of size TAILength
func TAIPack(t TAI) []byte {
	result := make([]byte, TAILength)
	binary.BigEndian.PutUint64(result[:], t.x)
	return result
}

// TAIUnpack unpacks a TAI timestamp from a byte array of size TAILength
func TAIUnpack(s []byte) TAI {
	return TAI{x: binary.BigEndian.Uint64(s[:])}
}

func (t TAI) String() string {
	buf := TAIPack(t)
	s := fmt.Sprintf("@%02X%02X%02X%02X%02X%02X%02X%02X",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6],
		buf[7])
	return s
}

// TAIfromString returns a TAI struct from an ASCII TAI representation
func TAIfromString(str string) (TAI, error) {
	if str[0] != '@' {
		return TAI{}, fmt.Errorf("TAI representation  %s is not valid, it does not begin with an '@'", str)
	}

	buf, err := hex.DecodeString(str[1:])
	if len(buf) != TAILength || err != nil {
		return TAI{}, err
	}

	return TAIUnpack(buf[:]), nil
}

// TAIfromTime returns a TAI struct from time.Time
func TAIfromTime(t time.Time) TAI {
	return TAI{x: TAICONST + lsoffset(t) + uint64(t.Unix())}
}
