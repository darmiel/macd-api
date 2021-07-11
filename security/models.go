package security

import "github.com/darmiel/macd-api/model"

var Accepted = []string{"NASDAQ", NYSE}

////////////////////////////////////////////////////

type NASDAQSecurity struct {
	Symbol          string `csv:"Symbol"`
	SecurityName    string `csv:"Security Name"`
	MarketCategory  string `csv:"Market Category"`
	TestIssue       bool   `csv:"Test Issue"`
	FinancialStatus string `csv:"Financial Status"`
	RoundLot        string `csv:"Round Lot Size"`
	ETF             bool   `csv:"ETF"`
	NextShares      bool   `csv:"NextShares"`
}

func (n *NASDAQSecurity) ToSymbol() *model.Symbol {
	return &model.Symbol{
		Symbol:   n.Symbol,
		Name:     n.SecurityName,
		ETF:      n.ETF,
		Exchange: "NASDAQ",
	}
}

////////////////////////////////////////////////////

type OtherSecurity struct {
	Symbol       string `csv:"ACT Symbol"`
	SecurityName string `csv:"Security Name"`
	Exchange     string `csv:"Exchange"`
	CQSSymbol    string `csv:"CQS Symbol"`
	ETF          bool   `csv:"ETF"`
	RoundLot     string `csv:"Round Lot Size"`
	TestIssue    bool   `csv:"Test Issue"`
	NASDAQSymbol string `csv:"NASDAQ Symbol"`
}

func (o *OtherSecurity) ToSymbol() *model.Symbol {
	return &model.Symbol{
		Symbol:   o.Symbol,
		Name:     o.SecurityName,
		ETF:      o.ETF,
		Exchange: o.Exchange,
	}
}
