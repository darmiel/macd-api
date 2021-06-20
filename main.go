package main

import (
	"fmt"
	"github.com/darmiel/macd-api/csv"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"strings"
	"time"
)

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

	fmt.Println("Retrieving file:", FTPFileNASDAQ)
	var resp *ftp.Response
	if resp, err = conn.Retr(FTPFileNASDAQ); err != nil {
		panic(err)
	}

	fmt.Println("Reading file ...")
	var buf []byte
	if buf, err = ioutil.ReadAll(resp); err != nil {
		panic(err)
	}

	fmt.Println("Parsing file ...")
	var parse *csv.CSVFile
	if parse, err = csv.Parse(buf, &csv.ParseOptions{
		Separator:  '|',
		CleanSpace: true,
		Blacklist:  []string{"File Creation Time: "},
	}); err != nil {
		panic(err)
	}
	fmt.Println("--------------------------------")

	var out []*NASDAQSecurity
	for i, v := range parse.Rows {
		dummy := new(NASDAQSecurity)
		if err = v.Unmarshal(dummy); err != nil {
			panic(err)
		}
		out = append(out, dummy)
		if dummy.ETF {
			fmt.Printf("%d :: %+v\n", i, dummy)
		}
	}
	fmt.Println(strings.Join(parse.Headers, " | "))
	fmt.Println("Parsed:", len(out), "securities.")
}
