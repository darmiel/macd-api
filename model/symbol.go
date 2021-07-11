package model

import "regexp"

var (
	SecurityPattern = regexp.MustCompile("^[A-Z]{1,5}$")
)

//

type Symbol struct {
	Symbol   string `gorm:"primaryKey"`
	Name     string
	ETF      bool
	Exchange string
}

type SymbolParser interface {
	ToSymbol() *Symbol
}

//

func ConvertToGenericArray(inp []*Symbol) (res []interface{}) {
	res = make([]interface{}, len(inp))
	for i, v := range inp {
		res[i] = v
	}
	return
}

func (s *Symbol) IsSymbolValid() bool {
	return SecurityPattern.MatchString(s.Symbol)
}

func (s *Symbol) IsAccepted(exchanges ...string) bool {
	if s.ETF || !s.IsSymbolValid() {
		return false
	}
	for _, e := range exchanges {
		if e == s.Exchange {
			return true
		}
	}
	return false
}
