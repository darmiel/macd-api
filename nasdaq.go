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
	Symbol          string
	SecurityName    string
	MarketCategory  string
	TestIssue       bool
	FinancialStatus string
	RoundLot        string
	ETF             bool
	NextShares      bool
}
