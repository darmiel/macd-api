package main

const FTPFileOther = "Symboldirectory/otherlisted.txt"

// Exchange
const (
	NYSE     = "N"
	NYSEMKT  = "A"
	NYSEARCA = "P"
	BATS     = "Z"
	IEXG     = "V"
)

type OtherSecurity struct {
	ACTSymbol    string
	SecurityName string
	Exchange     string
	CQSSymbol    string
	ETF          bool
	RoundLot     string
	TestIssue    bool
	NASDAQSymbol string
}
