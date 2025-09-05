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

// TAIN struct to store TAIN timestamps
type TAIN struct {
	sec  uint64
	nano uint32
}

// TAINNow returns the current timestamp in TAIN struct
func TAINNow() TAIN {
	now := time.Now()
	return TAIN{
		sec:  TAICONST + lsoffset(now) + uint64(now.Unix()),
		nano: uint32(now.Nanosecond()),
	}
}

// TAINAdd adds a time.Duration to a TAIN timestamp
func TAINAdd(a TAIN, b time.Duration) TAIN {
	var result TAIN
	result.sec = a.sec + uint64(b.Seconds())
	result.nano = a.nano + uint32(b.Nanoseconds()-int64(b.Seconds())*1000000000)
	if result.nano > 999999999 {
		result.sec++
		result.nano -= 1000000000
	}
	return result
}

// TAINSub subtracts two TAI timestamps
func TAINSub(a, b TAIN) (time.Duration, error) {
	s := a.sec - b.sec
	n := a.nano - b.nano
	if n > a.nano {
		n += 1000000000
		s--
	}
	q, err := time.ParseDuration(fmt.Sprintf("%ds%dns", s, n))
	return q, err
}

// TAINTime returns a go time object from a TAIN timestamp
func TAINTime(t TAIN) time.Time {
	tm := time.Unix(int64(t.sec-TAICONST), int64(t.nano)).UTC()
	return tm.Add(-time.Duration(lsoffset(tm)) * time.Second)
}

// TAINPack packs a TAIN timestamp in a byte array of size TAINLength
func TAINPack(t TAIN) []byte {
	result := make([]byte, TAINLength)
	binary.BigEndian.PutUint64(result[:], t.sec)
	binary.BigEndian.PutUint32(result[TAILength:], t.nano)
	return result
}

// TAINUnpack unpacks a TAIN timestamp from a byte array of size TAINLength
func TAINUnpack(s []byte) TAIN {
	var result TAIN
	result.sec = binary.BigEndian.Uint64(s[:])
	result.nano = binary.BigEndian.Uint32(s[TAILength:])
	return result
}

func (t TAIN) String() string {
	buf := TAINPack(t)
	s := fmt.Sprintf("@%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6],
		buf[7], buf[8], buf[9], buf[10], buf[11])
	return s
}

// TAINfromString returns a TAIN struct from an ASCII TAIN representation
func TAINfromString(str string) (TAIN, error) {
	if str[0] != '@' {
		return TAIN{}, fmt.Errorf("TAI representation  %s is not valid, it does not begin with an '@'", str)
	}

	buf, err := hex.DecodeString(str[1:])
	if len(buf) != TAINLength || err != nil {
		return TAIN{}, err
	}

	return TAINUnpack(buf[:]), nil
}

// TAINfromTime returns a TAIN struct from time.Time
func TAINfromTime(t time.Time) TAIN {
	return TAIN{
		sec:  TAICONST + lsoffset(t) + uint64(t.Unix()),
		nano: uint32(t.Nanosecond())}
}
