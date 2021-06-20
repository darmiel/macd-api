package csv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CSVRow struct {
	Headers *[]string
	Values  []string
}

type CSVFile struct {
	Headers   []string
	Rows      []*CSVRow
	Separator string
}

/////////////// UTIL

func escapedRuneSplit(str string, sep rune) (resp []string) {
	if str == "" {
		return []string{}
	}
	str = str + string(sep)
	var capt strings.Builder
	for i, r := range str {
		if r == sep {
			if i == 0 || str[i-1] != '\\' {
				resp = append(resp, capt.String())
				capt.Reset()
				continue
			}
		} else if r == '\\' && i+1 < len(str) && rune(str[i+1]) == sep {
			continue
		}
		capt.WriteRune(r)
	}
	return
}

/////////////// OPTS

type ParseOptions struct {
	Separator   rune
	HeaderIndex int
	CleanSpace  bool
	Blacklist   []string
}

func defaultOpts(options ...*ParseOptions) (resp *ParseOptions) {
	resp = &ParseOptions{
		Separator:   ',',
		HeaderIndex: 0,
		CleanSpace:  false,
	}
	// Fill defaults
	for _, opt := range options {
		if opt.Separator != rune(0) {
			resp.Separator = opt.Separator
		}
		if opt.HeaderIndex != 0 {
			resp.HeaderIndex = opt.HeaderIndex
		}
		if opt.CleanSpace != false {
			resp.CleanSpace = opt.CleanSpace
		}
		if len(opt.Blacklist) > 0 {
			resp.Blacklist = append(resp.Blacklist, opt.Blacklist...)
		}
	}
	return
}

/////////////// PARSE

func Parse(buf []byte, options ...*ParseOptions) (res *CSVFile, err error) {
	opt := defaultOpts(options...)
	chkbl := len(opt.Blacklist) > 0

	res = new(CSVFile)
	var hlen int

	for i, line := range strings.Split(string(buf), "\n") {
		// skip empty lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// data
		data := escapedRuneSplit(line, opt.Separator)
		if opt.CleanSpace {
			for j, v := range data {
				data[j] = strings.TrimSpace(v)
			}
		}

		// header?
		if i == opt.HeaderIndex {
			res.Headers = data
			hlen = len(res.Headers)
			continue
		}

		// row
		if len(data) != hlen {
			// invalid line
			// raise error?
			fmt.Println("ERR: Line", i+1, "has an invalid amount of data")
			continue
		}

		// Blacklist checking
		if chkbl {
			skip := false
			for _, d := range data {
				for _, b := range opt.Blacklist {
					if strings.Contains(d, b) {
						skip = true
						break
					}
				}
			}
			if skip {
				continue
			}
		}

		res.Rows = append(res.Rows, &CSVRow{Headers: &res.Headers, Values: data})
	}

	return
}

func toBool(str string) bool {
	switch strings.ToLower(str) {
	case "y", "yes", "1", "true":
		return true
	}
	return false
}

func (r *CSVRow) Unmarshal(target interface{}) (err error) {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("invalid type: %T", target)
	}
	if !rv.IsValid() {
		return fmt.Errorf("target is invalid")
	}

	vo := reflect.ValueOf(target).Elem()
	to := vo.Type()

	num := vo.NumField()
	for i := 0; i < num; i++ {
		tf := to.Field(i)
		vf := vo.Field(i)
		exp := tf.Tag.Get("csv") // tag

		// skip if no csv tag attached
		if exp == "" || exp == "-" {
			continue
		}

		var index = -1
		for j, h := range *r.Headers {
			if strings.EqualFold(h, exp) {
				index = j
				break
			}
		}
		// field not found
		if index < 0 {
			continue
		}

		// get value
		value := r.Values[index]

		switch vf.Kind() {
		case reflect.String:
			vf.SetString(value)
		case reflect.Bool:
			bf := toBool(value)
			vf.SetBool(bf)
		case reflect.Int, reflect.Int64:
			var iv int
			if iv, err = strconv.Atoi(value); err != nil {
				return
			}
			vf.SetInt(int64(iv))
		}
	}
	return
}
