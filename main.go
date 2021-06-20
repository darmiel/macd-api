package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"strings"
	"time"
)

// ftp.nasdaqtrader.com/Symboldirectory/nasdaqlisted.txt
// ftp.nasdaqtrader.com/Symboldirectory/otherlisted.txt
//
const FTPFile = "Symboldirectory/nasdaqlisted.txt"

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

func main() {
	fmt.Println("Connecting to FTP ...")
	conn, err := ftp.Dial("ftp.nasdaqtrader.com:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	fmt.Println("Logging in with anonymous ...")
	if err = conn.Login("anonymous", "anonymous"); err != nil {
		panic(err)
	}

	fmt.Println("Retrieving file:", FTPFile)
	var resp *ftp.Response
	if resp, err = conn.Retr(FTPFile); err != nil {
		panic(err)
	}

	fmt.Println("Reading file ...")
	var buf []byte
	if buf, err = ioutil.ReadAll(resp); err != nil {
		panic(err)
	}

	fmt.Println("Parsing file ...")
	for i, line := range strings.Split(string(buf), "\n") {
		// skip header
		if i == 0 {
			continue
		}
		spl := strings.Split(line, "|")
		if len(spl) != 8 {
			fmt.Println("ERR", line)
			continue
		}

		sec := &NASDAQSecurity{
			Symbol:          spl[0],
			SecurityName:    spl[1],
			MarketCategory:  spl[2],
			TestIssue:       spl[3] == "Y",
			FinancialStatus: spl[4],
			RoundLot:        spl[5],
			ETF:             spl[6] == "Y",
			NextShares:      spl[7] == "Y",
		}
		fmt.Printf("+ %+v\n", sec)
	}
}
