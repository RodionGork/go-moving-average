package counter

import (
	"math"
	"runtime"
	"testing"
)

func TestCountFiveMinutes(t *testing.T) {
	var timeNow = int64(1000)
	timeNowUnixFunc = func() int64 {
		return timeNow
	}
	cnt := createCounter()
	for i := 0; i < 300; i++ {
		cnt.ch <- int64(200 + i)
		timeNow += 1
		runtime.Gosched() // this is not 100% robust way to ensure cnt receiver works meanwhile, but let it be
	}
	avg, tot := cnt.averageAndCount()
	if math.Abs(float64(tot)-60.0) > 5 || math.Abs(float64(avg)-(200+300-60/2)) > 5 {
		t.Errorf("Unexpected results: avg=%d total=%d", avg, tot)
	}
}
