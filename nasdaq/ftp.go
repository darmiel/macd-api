package nasdaq

import (
	"fmt"
	"github.com/darmiel/macd-api/csv"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"time"
)

func FetchFile(rfile string) (buf []byte, err error) {
	var conn *ftp.ServerConn
	if conn, err = ftp.Dial("ftp.nasdaqtrader.com:21", ftp.DialWithTimeout(5*time.Second)); err != nil {
		return
	}
	if err = conn.Login("anonymous", "anonymous"); err != nil {
		return
	}
	var resp *ftp.Response
	if resp, err = conn.Retr(rfile); err != nil {
		return
	}
	buf, err = ioutil.ReadAll(resp)
	return
}

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
	return

}

/*
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

func requestPack(id int, pack []*NASDAQSecurity, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(pack); i++ {
		s := pack[i]

		url := fmt.Sprintf("", s.Symbol)

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

		fmt.Println("OKE | worker", id, "go response for symbol", s.Symbol, "::",
			resp.Response().StatusCode, len(to), "historical entries")
	}
}

*/
