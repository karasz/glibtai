package glibtai

import (
	"testing"
	"time"
)

func TestLsoffset(t *testing.T) {
	mp := make(map[int]uint64)
	mp[1933] = 0
	mp[1982] = 21
	mp[2933] = lsoffset(time.Now())
	for i, q := range mp {
		x := time.Date(i, time.August, 1, 0, 0, 0, 0, time.UTC)
		z := lsoffset(x)
		if z != q {
			t.Errorf("Offset for %v should be %d, not %d", x, q, z)
		}
	}
}
