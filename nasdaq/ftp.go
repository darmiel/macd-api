package nasdaq

import (
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

func FetchAll() (out []SecurityModel, err error) {
	var (
		nd    []*NASDAQSecurity
		ot    []*OtherSecurity
		index int
	)

	if nd, err = FetchNASDAQ(); err != nil {
		return
	}
	if ot, err = FetchOther(); err != nil {
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
