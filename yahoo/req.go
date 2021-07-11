package yahoo

import (
	"errors"
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/model"
	"github.com/imroc/req"
	"sort"
	"time"
)

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrHighEmpty       = errors.New("highs empty")
	ErrLowEmpty        = errors.New("lows empty")
	ErrVolumeEmpty     = errors.New("volumes empty")
	ErrCloseEmpty      = errors.New("closes empty")
	ErrOpenEmpty       = errors.New("opens empty")
)

func Historical90From(data []*model.Historic) (res model.Quarter, err error) {
	if len(data) < 90 {
		return model.Quarter{}, fmt.Errorf("required 90 records but %d provided", len(data))
	}
	r := data[:] // copy data
	sort.Slice(r, func(i, j int) bool {
		return r[i].DayDate.Before(r[j].DayDate)
	})
	if len(r) > 90 {
		r = r[len(r)-90:]
	}
	// fill
	for i, v := range r {
		res[i] = v
	}
	return
}

func RequestHistorical(symbol, interval, rng string) (resp []*model.Historic, err error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?formatted=true&interval=%s&range=%s",
		symbol, interval, rng)

	var http *req.Resp

	// try GET request 3 times
	for i := 0; i < 3; i++ {
		if http, err = req.Get(url); err != nil {
			if http != nil {
				fmt.Println("WARN :: Invalid response:", http.String(), "| try", i+1)
			} else {
				fmt.Println("WARN :: Invalid response ¯\\_(ツ)_/¯ | try", i+1)
			}
			// sleep for 1 sec
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if http == nil {
		return nil, errors.New("`http` was empty")
	}

	hr := new(historicalResponse)
	if err = http.ToJSON(hr); err != nil {
		return
	}

	resp, err = hr.To()
	return
}

type historicalResponse struct {
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

func (inp *historicalResponse) To() (out []*model.Historic, err error) {
	if inp.Chart.Result == nil {
		return nil, ErrInvalidResponse
	}
	if inp.Chart.Error != nil {
		return nil, fmt.Errorf("%v", inp.Chart.Error)
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
	for _, res := range inp.Chart.Result {
		if res.Meta.Symbol != "" {
			symbol = res.Meta.Symbol
		}
		if res.Timestamp != nil {
			timestamps = res.Timestamp
		}
		if quotes := res.Indicators.Quote; quotes != nil {
			for _, qv := range quotes {
				if qv.Close != nil {
					closes = qv.Close
				}
				if qv.High != nil {
					highs = qv.High
				}
				if qv.Low != nil {
					lows = qv.Low
				}
				if qv.Open != nil {
					opens = qv.Open
				}
				if qv.Volume != nil {
					volumes = qv.Volume
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
		high := highs[i]
		origT := time.Unix(t, 0)
		out = append(out, &model.Historic{
			Symbol:   symbol,
			DayDate:  common.NormalizeTimeNoon(origT),
			OrigDate: origT,
			High:     high,
			Low:      lows[i],
			Volume:   volumes[i],
			Close:    closes[i],
			Open:     opens[i],
		})
	}

	return
}
