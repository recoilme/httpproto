package httpproto

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

var (
	crlfcrlf = []byte("\r\n\r\n")
	crlf     = []byte("\r\n")
	space    = []byte(" ")
	postbin  = []byte("POST")

	contentlen = []byte("Content-Length")
)

// Request-Line   = Method SP Request-URI SP HTTP-Version CRLF
// []headers\r\n\r\n
// body

// Parse - read headers, parse content-length, and return parsed request as []byte and leftover, or nil & leftover
// return - leftover, request
func Parse(b []byte) ([]byte, []byte) {
	if len(b) == 0 { // that's not java - it's safe for nil
		return nil, nil
	}
	if i := bytes.Index(b, crlfcrlf); i >= 0 {
		if i == 0 {
			//if start from crlfcrlf - read crlfcrlf
			return b[i+len(crlfcrlf):], nil
		}

		headers := bytes.Split(b[:i+len(crlfcrlf)], crlf)
		cntlen := 0
		for _, header := range headers {
			//fmt.Printf("header:%+v\n", string(header))
			if bytes.HasPrefix(header, contentlen) {
				fields := bytes.Split(header, space)
				len, err := strconv.Atoi(string(fields[len(fields)-1]))
				if err == nil && len > 0 {
					cntlen = len
					break
				}
			}
		}
		//println("l1", len(b), "l2", (i + len(crlfcrlf) + cntlen))
		//fmt.Printf("%+v''\n", (string(b[i+len(crlfcrlf):])))
		if len(b) < (i + len(crlfcrlf) + cntlen) {
			return b, nil
		}
		return b[(i + len(crlfcrlf) + cntlen):], b[:(i + len(crlfcrlf) + cntlen)]

	}
	return b, nil
}

func scanRequestLine(line []byte) (method, uri, version string, err error) {
	pattern := "%s %s %s\r\n"
	dest := []interface{}{&method, &uri, &version}
	n, err := fmt.Sscanf(string(line), pattern, dest...)
	if n != len(dest) {
		err = errors.New("scanRequestLine error:" + string(line))
	}
	return
}
