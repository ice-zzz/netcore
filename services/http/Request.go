/**
 *     ______                 __
 *    /\__  _\               /\ \
 *    \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *       \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *        \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *        /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *        \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 *
 *                                                                    @寒冰
 *                                                            www.icezzz.cn
 *                                                     hanbin020706@163.com
 */
package http

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
)

type Request struct {
	Proto      string
	Method     string
	Path       string
	Query      string
	RemoteAddr string
	Connection string
	Header     map[string]string
}

func (req *Request) GetHeader(key string) string {
	return req.Header[key]
}

func (req *Request) Parsereq(data []byte) (n int, out []byte, err error) {
	sdata := string(data)
	var i, s int
	for k := range req.Header {
		delete(req.Header, k)
	}
	var line string
	var clen int
	var q = -1
	// method, path, proto line
	req.Proto = ""
	i = bytes.IndexByte(data, 32)
	if i == -1 {
		return
	}
	req.Method = sdata[:i]
	l := len(sdata)
	for i, s = i+1, i+1; i < l; i++ {
		if data[i] == 63 && q == -1 {
			q = i
		} else if data[i] == 32 {
			if q != -1 {
				req.Path = sdata[s:q]
				req.Query = sdata[q+1 : i]
			} else {
				req.Path = sdata[s:i]
			}
			i++
			s = bytes.Index(data[i:], []byte{13, 10})
			if s > -1 {
				s += i
				req.Proto = sdata[i:s]
			}
			break
		}
	}
	switch req.Proto {
	case "HTTP/1.0":
		req.Connection = "close"
	case "HTTP/1.1":
		req.Connection = "keep-alive"
	default:
		return 0, nil, fmt.Errorf("malformed request")
	}
	for s += 2; s < l; s += i + 2 {
		i = bytes.Index(data[s:], []byte{13, 10})
		line = sdata[s : s+i]
		if i > 15 {
			switch {
			case line[:15] == "Content-Length:", line[:15] == "Content-length:":
				clen, _ = strconv.Atoi(line[16:])
			case line == "Connection: close", line == "Connection: Close":
				req.Connection = "close"
			default:
				j := bytes.IndexByte(data[s:s+i], 58)
				req.Header[line[:j]] = line[j+2:]
			}
		} else if i == 0 {
			s += i + 2
			if clen == 0 && req.Header["Transfer-Encoding"] == "chunked" {

				for ; s < l; s += 2 {
					i = bytes.Index(data[s:], []byte{13, 10})
					if i == -1 {
						return 0, nil, nil
					}
					b := make([]byte, 8)
					if i&1 == 0 {
						_, _ = hex.Decode(b[8-i/2:], data[s:s+i])
					} else {
						tmp, _ := hex.DecodeString("0" + sdata[s:s+i])
						copy(b[7-i/2:], tmp)

					}
					clen = int(b[0])<<56 | int(b[1])<<48 | int(b[2])<<40 | int(b[3])<<32 | int(b[4])<<24 | int(b[5])<<16 | int(b[6])<<8 | int(b[7])
					s += i + 2
					if l-s < clen {
						return 0, nil, nil
					}
					if clen > 0 {
						out = append(out, data[s:s+clen]...)
						s += clen
					} else if l-s == 2 && data[s] == 13 && data[s+1] == 10 {
						return s + 2, out, nil
					}

				}

			} else {
				if l-s < clen {
					return 0, nil, nil
				}
				return s + clen, data[s : s+clen], nil
			}
		} else {
			j := bytes.IndexByte(data[s:s+i], 58)
			req.Header[line[:j]] = line[j+2:]
		}

	}

	// not enough data
	return 0, nil, nil
}
