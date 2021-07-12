package security

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/csv"
	"github.com/darmiel/macd-api/model"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"time"
)

type ftpConn struct {
	*ftp.ServerConn
}

func FTPConn() (conn *ftpConn, err error) {
	var serverConn *ftp.ServerConn
	if serverConn, err = ftp.Dial("ftp.nasdaqtrader.com:21", ftp.DialWithTimeout(10*time.Second)); err != nil {
		return
	}
	if err = serverConn.Login("anonymous", "anonymous"); err != nil {
		return
	}
	return &ftpConn{serverConn}, nil
}

func MustFTPConn() *ftpConn {
	conn, err := FTPConn()
	if err != nil {
		panic(err)
	}
	return conn
}

func (conn *ftpConn) FetchFile(rfile string) (buf []byte, err error) {
	var resp *ftp.Response
	if resp, err = conn.Retr(rfile); err != nil {
		fmt.Println("errrr:", err)
		return
	}
	buf, err = ioutil.ReadAll(resp)
	return
}

func FetchAll() (out []*model.Symbol, err error) {
	var (
		nd    []*NASDAQSecurity
		ot    []*OtherSecurity
		index int
	)

	fmt.Println(common.Info(), "Fetching NASDAQ ...")
	if nd, err = MustFTPConn().FetchNASDAQ(); err != nil {
		fmt.Println(common.Error(), "Error:", err)
		return
	}
	fmt.Println(common.Info(), "Fetching Other ...")
	if ot, err = MustFTPConn().FetchOther(); err != nil {
		fmt.Println(common.Error(), "Error:", err)
		return
	}

	out = make([]*model.Symbol, len(nd)+len(ot))

	for _, n := range nd {
		out[index] = n.ToSymbol()
		index++
	}
	for _, o := range ot {
		out[index] = o.ToSymbol()
		index++
	}

	for _, o := range out {
		o.Use = o.IsAccepted(Accepted...)
	}

	return
}

func FetchAllAccepted() (out []*model.Symbol, err error) {
	var all []*model.Symbol
	if all, err = FetchAll(); err != nil {
		return
	}
	for _, a := range all {
		if a.IsAccepted(Accepted...) {
			out = append(out, a)
		}
	}
	return
}

func FetchAllValid() (out []*model.Symbol, err error) {
	var all []*model.Symbol
	if all, err = FetchAll(); err != nil {
		return
	}
	for _, a := range all {
		if a.IsSymbolValid() {
			out = append(out, a)
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////
// NASDAQ

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

func (conn *ftpConn) FetchNASDAQ() (out []*NASDAQSecurity, err error) {
	var buf []byte
	if buf, err = conn.FetchFile(FTPFileNASDAQ); err != nil {
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

///////////////////////////////////////////////////////////////////////////
// NYSE

const FTPFileOther = "Symboldirectory/otherlisted.txt"

func (conn *ftpConn) FetchOther() (out []*OtherSecurity, err error) {
	var buf []byte
	if buf, err = conn.FetchFile(FTPFileOther); err != nil {
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

const (
	NASDAQ = "NASDAQ"
)
