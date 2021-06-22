package nasdaq

import "regexp"

var SecurityPattern = regexp.MustCompile("^[A-Z]{1,5}$")

type SecurityModel interface {
	Symbol() string
	SecurityName() string
	IsETF() bool
	Exchange() string
}

func IsSymbolValid(m SecurityModel) bool {
	return SecurityPattern.MatchString(m.Symbol())
}

////////////////////////////////////////////////////

type NASDAQSecurity struct {
	RawSymbol       string `csv:"Symbol"`
	RawSecurityName string `csv:"Security Name"`
	MarketCategory  string `csv:"Market Category"`
	TestIssue       bool   `csv:"Test Issue"`
	FinancialStatus string `csv:"Financial Status"`
	RoundLot        string `csv:"Round Lot Size"`
	ETF             bool   `csv:"ETF"`
	NextShares      bool   `csv:"NextShares"`
}

func (n *NASDAQSecurity) Symbol() string {
	return n.RawSymbol
}
func (n *NASDAQSecurity) SecurityName() string {
	return n.RawSecurityName
}
func (n *NASDAQSecurity) IsETF() bool {
	return n.ETF
}
func (n *NASDAQSecurity) Exchange() string {
	return "NASDAQ"
}

////////////////////////////////////////////////////

type OtherSecurity struct {
	RawACTSymbol    string `csv:"ACT Symbol"`
	RawSecurityName string `csv:"Security Name"`
	RawExchange     string `csv:"Exchange"`
	CQSSymbol       string `csv:"CQS Symbol"`
	ETF             bool   `csv:"ETF"`
	RoundLot        string `csv:"Round Lot Size"`
	TestIssue       bool   `csv:"Test Issue"`
	NASDAQSymbol    string `csv:"NASDAQ Symbol"`
}

func (o *OtherSecurity) Symbol() string {
	return o.RawACTSymbol
}
func (o *OtherSecurity) SecurityName() string {
	return o.RawSecurityName
}
func (o *OtherSecurity) IsETF() bool {
	return o.ETF
}
func (o *OtherSecurity) Exchange() string {
	return o.RawExchange
}
