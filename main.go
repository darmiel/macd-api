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
	for i, v := range parse.Values {
		fmt.Println(i, "=", strings.Join(v, " :: "))
	}
	fmt.Println(strings.Join(parse.Headers, " | "))
}
