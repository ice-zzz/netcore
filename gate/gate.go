package gate

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/ice-zzz/netcore/network"
)

type Gate struct {
	MaxConnNum int
	MaxMsgLen  uint32
	Processor  network.Processor

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration

	// tcp
	TCPAddr        string
	LenMsgLen      int
	LittleEndian   bool
	LenMsgLenInMsg bool
	ChanStop       bool

	wg sync.WaitGroup
}

func (gate *Gate) OnInit() {
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer

	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout

		wsServer.NewAgent = func(conn network.Conn) network.Agent {
			a := &agent{}
			a.conn = conn
			a.gate = gate
			a.userData = nil

			gate.wg.Add(1)
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.LenMsgLenInMsg = gate.LenMsgLenInMsg
		tcpServer.ChanStop = gate.ChanStop
		tcpServer.NewAgent = func(conn network.Conn) network.Agent {
			a := &agent{}
			a.conn = conn
			a.gate = gate
			a.userData = nil

			gate.wg.Add(1)
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
	gate.wg.Wait()
}

func (gate *Gate) OnDestroy() {

}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) React(b []byte) {

	if a.gate.Processor != nil {
		b, err := a.gate.Processor.Unmarshal(b)
		if err != nil {
			log.Printf("unmarshal message error: %v", err)
			return
		}
		err = a.gate.Processor.Route(b, a)
		if err != nil {
			log.Printf("route message error: %v", err)
			return
		}
	} else {
		log.Printf("agent not have a Processor")
		a.conn.Close()
	}

}

func (a *agent) OnClose() {

	a.gate.wg.Done()
}

func (a *agent) WriteMsg(msg []byte) {
	if a.gate.Processor != nil {
		b, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Printf("marshal message  error: %v", err)
		}
		a.conn.WriteMsg(b)
	} else {
		a.conn.WriteMsg(msg)
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
func (a *agent) OnInit() {

}
