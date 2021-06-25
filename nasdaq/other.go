package nasdaq

import (
	"github.com/darmiel/macd-api/csv"
	"github.com/jlaffaye/ftp"
)

const FTPFileOther = "Symboldirectory/otherlisted.txt"

func FetchOther(conn *ftp.ServerConn) (out []*OtherSecurity, err error) {
	var buf []byte
	if buf, err = FetchFile(conn, FTPFileOther); err != nil {
		return
	}

	var parse *csv.CSVFile
	if parse, err = csv.Parse(buf, &csv.ParseOptions{
		Separator:  '|',
		CleanSpace: true,
		Blacklist:  []string{"File Creation Time:"},
	}); err != nil {
		return
	}

	for _, v := range parse.Rows {
		dummy := new(OtherSecurity)
		if err = v.Unmarshal(dummy); err != nil {
			panic(err)
		}
		out = append(out, dummy)
	}
	return
}

// Exchange
const (
	NYSE     = "N"
	NYSEMKT  = "A"
	NYSEARCA = "P"
	BATS     = "Z"
	IEXG     = "V"
)
