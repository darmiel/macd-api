package main

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/csv"
	"github.com/darmiel/macd-api/keys"
	"github.com/imroc/req"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"strings"
	"sync"
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

	if len(out) == 0 {
		return
	}

	id := 0
	var wg sync.WaitGroup
	for _, idx := range common.Slice(len(out), 5) {
		var pack []*NASDAQSecurity
		for _, v := range idx {
			pack = append(pack, out[v])
		}
		wg.Add(1)
		go requestPack(id, pack, &wg)
		id++
	}

	wg.Wait()
	fmt.Println("Goroutines done!")
}

func requestPack(id int, pack []*NASDAQSecurity, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Requesting:", pack)

	for i := 0; i < len(pack); i++ {
		s := pack[i]

		// request api key
		var (
			key keys.APIKey
			err error
		)
		for {
			if key, err = keys.FindFreeKey(); err != nil {
				if err == keys.ErrNoFreeKey {
					// fmt.Println("WAT | worker", id, "Waiting for free key for symbol", s.Symbol, "...")
					time.Sleep(time.Second)
					continue
				}
				panic(err)
			}
			break
		}

		url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s",
			s.Symbol, key)

		fmt.Println("REQ requesting", url, "(Key:", key+")")
		resp, err := req.Get(url)
		if err != nil {
			fmt.Println("ERR | worker", id, "encountered error:", err, "on symbol", s.Symbol)
			continue
		}

		fmt.Println("OKE | worker", id, "go response for symbol", s.Symbol, "k:", key, "::",
			resp.Response().StatusCode, common.Smallify(resp.String(), 128))

		if strings.Contains(resp.String(), "Thank you for using") {
			fmt.Println("ERR | Key limit exceeded.")
			key.Invalidate()
			i--
		}
	}
}
