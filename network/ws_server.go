package network

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ice-zzz/netcore/ds/codec"
	"github.com/panjf2000/gnet"
)

type WSServer struct {
	Addr        string
	MaxConnNum  int
	MaxMsgLen   uint32
	CertFile    string
	KeyFile     string
	CaFile      string
	HTTPTimeout time.Duration
	NewAgent    func(Conn) Agent
	wg          sync.WaitGroup
	ConnNum     int32
	gnet.EventHandler
	Close func()
}

func (server *WSServer) Start() {

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		log.Printf("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}

	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 4096
		log.Printf("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		log.Printf("invalid HTTPTimeout, reset to %v", server.HTTPTimeout)
	}
	if server.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}

	option := []gnet.Option{gnet.WithTCPKeepAlive(time.Second * 600),
		gnet.WithCodec(&codec.WSCode{}),
		gnet.WithReusePort(true),
	}
	go gnet.Serve(server, "tcp://"+server.Addr, option...)
}
func (server *WSServer) OnInitComplete(svr gnet.Server) gnet.Action {
	// server.Close = svr.Close
	log.Printf("echo run websocks on " + server.Addr)
	return gnet.None
}
func (server *WSServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	num := int(atomic.AddInt32(&server.ConnNum, 1))
	if num >= server.MaxConnNum {
		log.Printf("too many connections")
		return nil, gnet.Close
	}

	return
}
func (server *WSServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	atomic.AddInt32(&server.ConnNum, -1)
	switch svr := c.Context().(type) {
	case *codec.WSconn:
		switch agent := svr.Ctx.(type) {
		case Agent:
			agent.OnClose()
		}
	}
	c.SetContext(nil)
	return
}
func (server *WSServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	switch svr := c.Context().(type) {
	case *codec.Httpserver:

		// tmp_ctx.Request.Body.Reset()
		// tmp_ctx.Request.Body.Write(data[len(data)-8 : len(data)])
		err := svr.Upgradews(c)
		if err != nil {
			action = gnet.Close
			return
		}
		if ws, ok := c.Context().(*codec.WSconn); ok {
			agent := server.NewAgent(gnetConn{c})
			agent.OnInit()
			ws.Ctx = agent
		} else {
			action = gnet.Close
			return
		}
		return
	case *codec.WSconn:
		if agent, ok := svr.Ctx.(Agent); ok {
			agent.React(frame)
		}
	}
	return
}
