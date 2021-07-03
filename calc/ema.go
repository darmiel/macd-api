package calc

import (
	"github.com/darmiel/macd-api/yahoo"
)

type T []float64

func (t T) Min(num int) float64 {
	l := len(t)
	if num > l {
		return -1
	}
	return t[l-1-num]
}

func (t T) Seven() (res [8]float64) {
	for i := 0; i < 8; i++ {
		res[i] = t.Min(i)
	}
	return
}

// requires data to be ordered by date ASC
// and min len of ({day} + 1)
func EMA(day int, data yahoo.Historical90) (res T, err error) {
	alpha := 2.0 / (float64(day) + 1.0)
	ema := float64(data[0].Close)
	for _, h := range data[1:] {
		ema = float64(h.Close)*alpha + ema*(1-alpha)
		res = append(res, ema)
	}
	return
}
