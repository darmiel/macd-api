package main

const FTPFileNASDAQ = "Symboldirectory/nasdaqlisted.txt"

// Market Category
const (
	NASDAQGlobalSelectMarketSM = "Q"
	NASDAQGlobalMarketSM       = "G"
	NASDAQCapitalMarket        = "S"
)

// Financial Status
const (
	Deficient                      = "D"
	Delinquent                     = "E"
	Bankrupt                       = "Q"
	Normal                         = "N"
	DeficientAndBankrupt           = "G"
	DeficientAndDelinquent         = "H"
	DelinquentAndBankrupt          = "J"
	DeficientDelinquentAndBankrupt = "K"
)

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
