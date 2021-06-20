package main

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/darmiel/macd-api/csv"
	"github.com/darmiel/macd-api/yahoo"
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
		if id >= 1 {
			break
		}

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

var (
	curNum = 0
	mu     sync.Mutex
)

func requestPack(id int, pack []*NASDAQSecurity, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(pack); i++ {
		s := pack[i]

		url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?formatted=true&interval=1d&range=1y", s.Symbol)

		mu.Lock()
		curNum++
		fmt.Println(curNum, "REQ requesting", url)
		mu.Unlock()

		resp, err := req.Get(url)
		if err != nil {
			fmt.Println("ERR | worker", id, "encountered error:", err, "on symbol", s.Symbol)
			continue
		}

		yr := new(yahoo.Response)
		if err := resp.ToJSON(yr); err != nil {
			fmt.Println("ERR | worker", id, "encountered error:", err, "on symbol", s.Symbol, "-- decoding")
			continue
		}

		to, err := yr.To()
		if err != nil {
			fmt.Println("ERR | worker", id, "encountered error:", err, "on symbol", s.Symbol, "-- wrapping")
			continue
		}

		for _, t := range to {
			fmt.Println("--")
			fmt.Printf("%+v\n", t)
			fmt.Println("---")
		}

		fmt.Println("OKE | worker", id, "go response for symbol", s.Symbol, "::",
			resp.Response().StatusCode, common.Smallify(resp.String(), 128))
	}
}
