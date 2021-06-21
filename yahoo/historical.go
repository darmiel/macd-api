package yahoo

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"gorm.io/gorm"
)

func RequestHistorical(symbol, interval, rng string) (resp []*Historical, err error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?formatted=true&interval=%s&range=%s",
		symbol, interval, rng)

	var http *req.Resp
	if http, err = req.Get(url); err != nil {
		if http != nil {
			fmt.Println("WARN :: Invalid response:", http.String())
		} else {
			fmt.Println("WARN :: Invalid response ¯\\_(ツ)_/¯")
		}
		return
	}

	hr := new(HistoricalResponse)
	if err = http.ToJSON(hr); err != nil {
		return
	}

	resp, err = hr.To()
	return
}

type HistoricalResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string `json:"currency"`
				Symbol               string `json:"symbol"`
				Exchang              string `json:"exchangeName"`
				FirstTradeDate       int    `json:"firstTradeDate"`
				TimeZone             string `json:"timezone"`
				Exchangetimezonename string `json:"exchangeTimezoneName"`
				DataGranularity      string `json:"dataGranularity"`
				Range                string `json:"range"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					High   []float32 `json:"high"`
					Low    []float32 `json:"low"`
					Volume []int     `json:"volume"`
					Close  []float32 `json:"close"`
					Open   []float32 `json:"open"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

type Historical struct {
	gorm.Model
	Symbol    string
	Timestamp int64
	High      float32
	Low       float32
	Volume    int
	Close     float32
	Open      float32
}

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrHighEmpty       = errors.New("highs empty")
	ErrLowEmpty        = errors.New("lows empty")
	ErrVolumeEmpty     = errors.New("volumes empty")
	ErrCloseEmpty      = errors.New("closes empty")
	ErrOpenEmpty       = errors.New("opens empty")
)

func (r *HistoricalResponse) To() (w []*Historical, err error) {
	if r.Chart.Result == nil {
		return nil, ErrInvalidResponse
	}
	if r.Chart.Error != nil {
		return nil, fmt.Errorf("%v", r.Chart.Error)
	}

	var (
		symbol     string
		timestamps []int64
		highs      []float32
		lows       []float32
		volumes    []int
		closes     []float32
		opens      []float32
	)

	//
	for _, res := range r.Chart.Result {
		if res.Meta.Symbol != "" {
			symbol = res.Meta.Symbol
		}
		if res.Timestamp != nil {
			timestamps = res.Timestamp
		}
		if q := res.Indicators.Quote; q != nil {
			for _, qu := range q {
				if qu.Close != nil {
					closes = qu.Close
				}
				if qu.High != nil {
					highs = qu.High
				}
				if qu.Low != nil {
					lows = qu.Low
				}
				if qu.Open != nil {
					opens = qu.Open
				}
				if qu.Volume != nil {
					volumes = qu.Volume
				}
			}
		}
	}

	if highs == nil {
		return nil, ErrHighEmpty
	}
	if lows == nil {
		return nil, ErrLowEmpty
	}
	if volumes == nil {
		return nil, ErrVolumeEmpty
	}
	if closes == nil {
		return nil, ErrCloseEmpty
	}
	if opens == nil {
		return nil, ErrOpenEmpty
	}

	for i, t := range timestamps {
		w = append(w, &Historical{
			Symbol:    symbol,
			Timestamp: t,
			High:      highs[i],
			Low:       lows[i],
			Volume:    volumes[i],
			Close:     closes[i],
			Open:      opens[i],
		})
	}

	return
}
