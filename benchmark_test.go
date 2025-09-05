// Copyright © 2018 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of glibtai, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>

package glibtai

import (
	"testing"
	"time"
)

func BenchmarkTAINow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TAINow()
	}
}

func BenchmarkTAIAdd(b *testing.B) {
	tai := TAINow()
	duration := time.Hour
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAIAdd(tai, duration)
	}
}

func BenchmarkTAISub(b *testing.B) {
	tai1 := TAINow()
	tai2 := TAIAdd(tai1, time.Hour)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TAISub(tai2, tai1)
	}
}

func BenchmarkTAITime(b *testing.B) {
	tai := TAINow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAITime(tai)
	}
}

func BenchmarkTAIPack(b *testing.B) {
	tai := TAINow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAIPack(tai)
	}
}

func BenchmarkTAIUnpack(b *testing.B) {
	tai := TAINow()
	packed := TAIPack(tai)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAIUnpack(packed)
	}
}

func BenchmarkTAIString(b *testing.B) {
	tai := TAINow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tai.String()
	}
}

func BenchmarkTAIfromString(b *testing.B) {
	tai := TAINow()
	str := tai.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TAIfromString(str)
	}
}

func BenchmarkTAIfromTime(b *testing.B) {
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAIfromTime(now)
	}
}

func BenchmarkTAINNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TAINNow()
	}
}

func BenchmarkTAINAdd(b *testing.B) {
	tain := TAINNow()
	duration := time.Hour
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAINAdd(tain, duration)
	}
}

func BenchmarkTAINSub(b *testing.B) {
	tain1 := TAINNow()
	tain2 := TAINAdd(tain1, time.Hour)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TAINSub(tain2, tain1)
	}
}

func BenchmarkTAINTime(b *testing.B) {
	tain := TAINNow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAINTime(tain)
	}
}

func BenchmarkTAINPack(b *testing.B) {
	tain := TAINNow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAINPack(tain)
	}
}

func BenchmarkTAINUnpack(b *testing.B) {
	tain := TAINNow()
	packed := TAINPack(tain)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAINUnpack(packed)
	}
}

func BenchmarkTAINString(b *testing.B) {
	tain := TAINNow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tain.String()
	}
}

func BenchmarkTAINfromString(b *testing.B) {
	tain := TAINNow()
	str := tain.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TAINfromString(str)
	}
}

func BenchmarkTAINfromTime(b *testing.B) {
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TAINfromTime(now)
	}
}

func BenchmarkLsoffset(b *testing.B) {
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lsoffset(now)
	}
}

func BenchmarkLsoffsetHistoric(b *testing.B) {
	historic := time.Date(1985, time.July, 1, 0, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lsoffset(historic)
	}
}

func BenchmarkLsoffsetPreLeap(b *testing.B) {
	preLeap := time.Date(1971, time.January, 1, 0, 0, 0, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lsoffset(preLeap)
	}
}
