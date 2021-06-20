package yahoo

import (
	"errors"
	"fmt"
)

type Response struct {
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
			Timestamp  []uint64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Volume []int     `json:"volume"`
					Close  []float64 `json:"close"`
					Open   []float64 `json:"open"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

type WrappedIndicators struct {
	Symbol    string
	Timestamp uint64
	High      float64
	Low       float64
	Volume    int
	Close     float64
	Open      float64
}

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrHighEmpty       = errors.New("highs empty")
	ErrLowEmpty        = errors.New("lows empty")
	ErrVolumeEmpty     = errors.New("volumes empty")
	ErrCloseEmpty      = errors.New("closes empty")
	ErrOpenEmpty       = errors.New("opens empty")
)

func (r *Response) To() (w []*WrappedIndicators, err error) {
	if r.Chart.Result == nil {
		return nil, ErrInvalidResponse
	}
	if r.Chart.Error != nil {
		return nil, fmt.Errorf("%v", r.Chart.Error)
	}

	var (
		symbol     string
		timestamps []uint64
		highs      []float64
		lows       []float64
		volumes    []int
		closes     []float64
		opens      []float64
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
		w = append(w, &WrappedIndicators{
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
