package calc

import (
	"github.com/darmiel/macd-api/model"
)

type T []float64

func (t T) Min(num int) float64 {
	l := len(t)
	if num > l {
		return -1
	}
	return t[l-1-num]
}

// T -> alias for Min
func (t T) T(num int) float64 {
	return t.Min(num)
}

func (t T) Seven() (res [8]float64) {
	for i := 0; i < 8; i++ {
		res[i] = t.Min(i)
	}
	return
}

// requires data to be ordered by date ASC
// and min len of ({day} + 1)
func EMA(day int, data model.Quarter) (res T) {
	alpha := 2.0 / (float64(day) + 1.0)
	ema := float64(data[0].Close)
	res = append(res, ema)
	for _, h := range data[1:] {
		ema = float64(h.Close)*alpha + ema*(1-alpha)
		res = append(res, ema)
	}
	return
}

func MACD(data model.Quarter) (res T) {
	ema10 := EMA(10, data)
	ema35 := EMA(35, data)
	for i := 0; i < len(data); i++ {
		res = append(res, ema10[i]-ema35[i])
	}
	return
}
