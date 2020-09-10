package network

import (
	"log"
	"net"

	"github.com/panjf2000/gnet"
)

type Conn interface {
	ReadMsg(buf []byte) ([]byte, error)
	WriteMsg(msg []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
}

// 新增一个适配Conn接口的结构
type gnetConn struct {
	c gnet.Conn
}

func (c gnetConn) ReadMsg(buf []byte) ([]byte, error) {
	log.Printf("gnet don't have ReadMsg interface,please use React")
	return nil, nil
}
func (c gnetConn) WriteMsg(msg []byte) error {
	c.c.AsyncWrite(msg)
	return nil
}
func (c gnetConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}
func (c gnetConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}
func (c gnetConn) Close() error {
	return c.c.Close()
}
