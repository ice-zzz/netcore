package codec

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/ice-zzz/netcore/ds"
	"github.com/klauspost/compress/gzip"
	"github.com/panjf2000/gnet"
)

type Httpserver struct {
	Request Request
	c       gnet.Conn
	Out     *ds.MsgBuffer
	Ws      *WSconn
	data    *bytes.Reader
}

type Request struct {
	Proto, Method string
	Path, Query   string
	RemoteAddr    string
	Connection    string
	Header        map[string]string
}

var Httppool = sync.Pool{New: func() interface{} {
	hs := &Httpserver{Out: new(ds.MsgBuffer)}
	hs.Request.Header = make(map[string]string)
	hs.data = &bytes.Reader{}
	return hs
}}

var msgbufpool = sync.Pool{New: func() interface{} {
	return new(ds.MsgBuffer)
}}
var gzippool = sync.Pool{New: func() interface{} {
	w, _ := gzip.NewWriterLevel(nil, 6)
	return w
}}

func (r *Request) GetHeader(key string) string {
	return r.Header[key]
}

func (hs *Httpserver) Ip(c gnet.Conn) (ip string) {

	if ip = hs.Request.GetHeader("X-Real-IP"); ip == "" {
		ip = c.RemoteAddr().String()
	}

	return ip
}
func (hs *Httpserver) IsMobile() bool {
	return false
}
func (hs *Httpserver) Lastvisit() int32 {
	return 0
}
func (hs *Httpserver) SetLastvisit(int32) {

}

func (hs *Httpserver) UserAgent() string {
	return hs.Request.GetHeader("UserAgent")
}

var errprotocol = errors.New("the client is not using the websocket protocol: ")

// http升级为websocket
func (hs *Httpserver) Upgradews(c gnet.Conn) (err error) {
	//
	hs.Out.Reset()
	/*if !(strings.Contains(c.Request.Head, "Connection: Upgrade")) {

		hs.Out.WriteString("HTTP/1.1 400 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")

		return errprotocol
	}*/
	if hs.Request.Method != "GET" {

		hs.Out.WriteString("HTTP/1.1 403 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")

		return errprotocol
	}
	/*
		if !(strings.Contains(c.Request.Head, "Sec-WebSocket-Extensions")) {

			hs.Out.WriteString("HTTP/1.1 400 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")

			return
		}*/

	/*if config.Server.Origin != "" && hs.Request.Header["Origin"] != config.Server.Origin {
		hs.Out.WriteString("HTTP/1.1 403 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")
		return errprotocol
	}*/
	if hs.Request.Header["Upgrade"] != "websocket" {
		hs.Out.WriteString("HTTP/1.1 403 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")

		return errprotocol
	}

	if hs.Request.Header["Sec-WebSocket-Version"] != "13" {
		hs.Out.WriteString("HTTP/1.1 403 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")

		return errprotocol
	}

	var challengeKey string

	if challengeKey = hs.Request.Header["Sec-WebSocket-Key"]; challengeKey == "" {
		hs.Out.WriteString("HTTP/1.1 403 Error\r\nContent-Type: text/plain\r\nContent-Length: 11\r\nConnection: close\r\n\r\nUnknonw MSG")

		return errprotocol
	}

	hs.Ws = &WSconn{
		IsServer:   true,
		ReadFinal:  true,
		Http:       hs,
		Write:      c.AsyncWrite,
		IsCompress: strings.Contains(hs.Request.Header["Sec-WebSocket-Extensions"], "permessage-deflate"),
		readbuf:    &ds.MsgBuffer{},
	}

	c.SetContext(hs.Ws)
	hs.Out.WriteString("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: ")
	hs.Out.WriteString(ComputeAcceptKey(challengeKey))
	hs.Out.WriteString("\r\n")
	if hs.Ws.IsCompress {
		hs.Out.WriteString("Sec-Websocket-Extensions: permessage-deflate; server_no_context_takeover; client_no_context_takeover\r\n")
	}
	hs.Out.WriteString("\r\n")
	hs.c.AsyncWrite(hs.Out.Bytes())
	return nil
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
						hex.Decode(b[8-i/2:], data[s:s+i])
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

var (
	static_patch = "./static"
	http1head200 = []byte("HTTP/1.1 200 OK\r\nserver: gnet by luyu6056\r\n")
	http1head206 = []byte("HTTP/1.1 206 Partial Content\r\nserver: gnet by luyu6056\r\n")
	http1head304 = []byte("HTTP/1.1 304 Not Modified\r\nserver: gnet by luyu6056\r\n")
	http1deflate = []byte("\r\nContent-encoding: deflate")
	http1gzip    = []byte("\r\nContent-encoding: gzip")
	http404b, _  = ioutil.ReadFile(static_patch + "/404.html")
	http1cache   = []byte("Cache-Control: max-age=86400\r\n")
	http1nocache = []byte("Cache-Control: no-store, no-cache, must-revalidate, max-age=0, s-maxage=0\r\nPragma: no-cache\r\n")
)

func (hs *Httpserver) Out404(err error) {
	hs.Out.WriteString("HTTP/1.1 404 Not Found\r\nContent-Length: ")
	hs.Out.WriteString(strconv.Itoa(len(http404b)))
	hs.Out.WriteString("\r\n\r\n")
	hs.Out.Write(http404b)
	hs.c.AsyncWrite(hs.Out.Bytes())
}

func init() {

	if http404b == nil {
		http404b = []byte("404 not found")
	}
}
