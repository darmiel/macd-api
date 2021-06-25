package nasdaq

import (
	"fmt"
	"github.com/darmiel/macd-api/common"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"time"
)

func Connection() (conn *ftp.ServerConn, err error) {
	if conn, err = ftp.Dial("ftp.nasdaqtrader.com:21", ftp.DialWithTimeout(10*time.Second)); err != nil {
		return
	}
	if err = conn.Login("anonymous", "anonymous"); err != nil {
		return
	}
	return
}

func MustConnection() *ftp.ServerConn {
	conn, err := Connection()
	if err != nil {
		panic(err)
	}
	return conn
}

func FetchFile(conn *ftp.ServerConn, rfile string) (buf []byte, err error) {
	var resp *ftp.Response
	if resp, err = conn.Retr(rfile); err != nil {
		fmt.Println("errrr:", err)
		return
	}
	buf, err = ioutil.ReadAll(resp)
	return
}

func FetchAll() (out []SecurityModel, err error) {
	var (
		nd    []*NASDAQSecurity
		ot    []*OtherSecurity
		index int
	)

	fmt.Println(common.Info(), "Fetching NASDAQ ...")
	if nd, err = FetchNASDAQ(MustConnection()); err != nil {
		fmt.Println(common.Error(), "Error:", err)
		return
	}
	fmt.Println(common.Info(), "Fetching Other ...")
	if ot, err = FetchOther(MustConnection()); err != nil {
		fmt.Println(common.Error(), "Error:", err)
		return
	}

	out = make([]SecurityModel, len(nd)+len(ot))

	for _, n := range nd {
		out[index] = n
		index++
	}
	for _, o := range ot {
		out[index] = o
		index++
	}

	return
}
func FetchAllAccepted() (out []SecurityModel, err error) {
	var all []SecurityModel
	if all, err = FetchAll(); err != nil {
		return
	}
	for _, a := range all {
		if IsModelAccepted(a) {
			out = append(out, a)
		}
	}
	return
}
