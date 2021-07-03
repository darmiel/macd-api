package calc

import (
	"github.com/darmiel/macd-api/yahoo"
)

// requires data to be ordered by date ASC
// and min len of ({day} + 1)
func EMA(day int, data []*yahoo.Historical) (res []float64, err error) {
	alpha := 2.0 / (float64(day) + 1.0)
	ema := float64(data[0].Close)
	i := 1
	for _, h := range data[1:] {
		i++
		ema = float64(h.Close)*alpha + ema*(1-alpha)
		if i > 90-8 {
			res = append(res, ema)
		}
	}
	return
}
