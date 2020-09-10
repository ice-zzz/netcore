package network

import (
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/panjf2000/gnet"
)

type TCPServer struct {
	connsNum   int32
	Addr       string
	MaxConnNum int
	NewAgent   func(Conn) Agent

	// msg parser
	LenMsgLen      int
	MinMsgLen      uint32
	MaxMsgLen      uint32
	LittleEndian   bool
	LenMsgLenInMsg bool // if the msg len contain in "header" len
	msgParser      *MsgParser
	ChanStop       bool
	gnet.EventHandler
	Close func()
}

func (server *TCPServer) Start() {
	server.init()
}

func (server *TCPServer) init() {

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		log.Printf("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}

	if server.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(server.LenMsgLen, server.MinMsgLen, server.MaxMsgLen)
	msgParser.SetByteOrder(server.LittleEndian)
	msgParser.SetLenMsgLenInMsg(server.LenMsgLenInMsg)

	go gnet.Serve(server, "tcp://"+server.Addr,
		gnet.WithTCPKeepAlive(time.Second*600),
		gnet.WithCodec(msgParser),
		gnet.WithReusePort(true),
	)
}
func (server *TCPServer) OnInitComplete(svr gnet.Server) gnet.Action {
	// server.Close = svr.Close
	log.Printf("echo run tcpserver on " + server.Addr)
	return gnet.None
}
func (server *TCPServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	num := int(atomic.AddInt32(&server.connsNum, 1))
	if num >= server.MaxConnNum {
		log.Printf("too many connections")
		return nil, gnet.Close
	}
	agent := server.NewAgent(gnetConn{c})
	c.SetContext(agent)
	agent.OnInit()
	return
}
func (server *TCPServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("gnet close %v", err)
	atomic.AddInt32(&server.connsNum, -1)
	switch agent := c.Context().(type) {
	case Agent:
		agent.OnClose()
	}
	c.SetContext(nil)
	return
}
func (server *TCPServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	defer Recover()
	switch agent := c.Context().(type) {
	case Agent:
		agent.React(frame)

		// default:
		// 	return gnet.Close
	}
	return
}
func Recover() {
	if err := recover(); err != nil {
		stack := debug.Stack()
		log.Printf(string(stack))
	}
}
