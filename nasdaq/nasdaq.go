package nasdaq

import (
	"github.com/darmiel/macd-api/csv"
)

const FTPFileNASDAQ = "Symboldirectory/nasdaqlisted.txt"

func FetchNASDAQ() (out []*NASDAQSecurity, err error) {
	var buf []byte
	if buf, err = FetchFile(FTPFileNASDAQ); err != nil {
		return
	}

	var parse *csv.CSVFile
	if parse, err = csv.Parse(buf, &csv.ParseOptions{
		Separator:  '|',
		CleanSpace: true,
		Blacklist:  []string{"File Creation Time: "},
	}); err != nil {
		return
	}

	for _, v := range parse.Rows {
		dummy := new(NASDAQSecurity)
		if err = v.Unmarshal(dummy); err != nil {
			panic(err)
		}
		out = append(out, dummy)
	}
	return
}

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
