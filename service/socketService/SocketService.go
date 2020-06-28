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
package socketService

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ice-zzz/netcore/internal/netpoll"
	"github.com/ice-zzz/netcore/manager/network"
	"github.com/ice-zzz/netcore/service"
	"github.com/ice-zzz/netcore/utils/gopool"
	"github.com/segmentio/ksuid"
)

type SocketServer struct {
	exit chan struct{}
	nm   *network.NetManager
	service.Entity
}

type deadliner struct {
	net.Conn
	t time.Duration
}

func New() *SocketServer {

	return &SocketServer{
		exit: make(chan struct{}),
		nm:   network.NewNetManager(),
		Entity: service.Entity{
			Name: "",
			Ip:   "0.0.0.0",
			Port: 5678,
		},
	}

}

func (serv *SocketServer) Start() {

	handle := func(conn net.Conn) {

		safeConn := deadliner{conn, time.Millisecond * 100}

		user := serv.nm.Group.Register(fmt.Sprintf("%s_%s", "s", ksuid.New().String()), safeConn)
		desc := netpoll.Must(netpoll.HandleRead(conn))

		fmt.Printf("用户 %s 进入 ip-> %s  \n", user.GetName(), conn.RemoteAddr().String())

		_ = serv.nm.Poller.Start(desc, func(ev netpoll.Event) {
			// 断线处理
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {

				_ = serv.nm.Poller.Stop(desc)
				serv.nm.Group.Remove(user)
				return
			}

			serv.nm.Pool.Schedule(func() {

				if err := user.Receive(); err != nil {

					_ = serv.nm.Poller.Stop(desc)
					serv.nm.Group.Remove(user)
				}
			})
		})

	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serv.Ip, serv.Port))
	if err != nil {
		log.Fatal(err)
		return
	}

	// fmt.Printf("websocket 正在监听端口-> %d \n", serv.Port)

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		ln, netpoll.EventRead|netpoll.EventOneShot,
	))

	accept := make(chan error, 1)

	_ = serv.nm.Poller.Start(acceptDesc, func(e netpoll.Event) {

		err := serv.nm.Pool.ScheduleTimeout(time.Millisecond, func() {
			conn, err := ln.Accept()
			if err != nil {
				accept <- err
				return
			}

			accept <- nil
			handle(conn)
		})
		if err == nil {
			err = <-accept
		}
		if err != nil {
			if err != gopool.ErrScheduleTimeout {
				goto cooldown
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				goto cooldown
			}

			fmt.Printf("连接错误: %v \n", err)

		cooldown:
			delay := 5 * time.Millisecond
			fmt.Printf("连接错误: %v; %s 秒后重试! \n ", err, delay)
			time.Sleep(delay)
		}

		_ = serv.nm.Poller.Resume(acceptDesc)
	})

	<-serv.exit

}

func (serv *SocketServer) Stop() {
	serv.exit <- struct{}{}
}

func (serv *SocketServer) AddHandler(messageType uint16, fun network.RecvHandler) {
	serv.nm.Group.Hanlder.AddHandler(messageType, fun)
}

func nameConn(conn net.Conn) string {
	return conn.LocalAddr().String() + " > " + conn.RemoteAddr().String()
}

func (d deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}
