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
package websocketService

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/ice-zzz/netcore/internal/netpoll"
	"github.com/ice-zzz/netcore/manager/network"
	"github.com/ice-zzz/netcore/utils/gopool"

	"github.com/ice-zzz/netcore/service"
)

type WebSocketServer struct {
	exit chan struct{}
	nm   *network.NetManager
	service.Entity
}

type deadliner struct {
	net.Conn
	t time.Duration
}

func New() *WebSocketServer {

	return &WebSocketServer{
		exit: make(chan struct{}),
		nm:   network.NewNetManager(),
		Entity: service.Entity{
			Name: "",
			Ip:   "0.0.0.0",
			Port: 5678,
		},
	}

}

func (webserv *WebSocketServer) Start() {

	handle := func(conn net.Conn) {

		safeConn := deadliner{conn, time.Millisecond * 100}

		_, err := ws.Upgrade(safeConn)
		if err != nil {
			fmt.Printf("%s: 升级失败: %v \n", nameConn(conn), err)
			_ = conn.Close()
			return
		}

		user := webserv.nm.Group.Register(safeConn)
		desc := netpoll.Must(netpoll.HandleRead(conn))

		fmt.Printf("用户 %s 进入 ip-> %s  \n", user.GetName(), conn.RemoteAddr().String())

		_ = webserv.nm.Poller.Start(desc, func(ev netpoll.Event) {
			// 断线处理
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {

				_ = webserv.nm.Poller.Stop(desc)
				webserv.nm.Group.Remove(user)
				return
			}

			webserv.nm.Pool.Schedule(func() {

				if err := user.Receive(); err != nil {

					_ = webserv.nm.Poller.Stop(desc)
					webserv.nm.Group.Remove(user)
				}
			})
		})

	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", webserv.Ip, webserv.Port))
	if err != nil {
		log.Fatal(err)
		return
	}

	// fmt.Printf("websocket 正在监听端口-> %d \n", webserv.Port)

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		ln, netpoll.EventRead|netpoll.EventOneShot,
	))
	// webserv.EventDispatch.DispatchEvent(COMPLETE, &events.Event{
	// 	EventType: COMPLETE,
	// 	Data:      nil,
	// })

	accept := make(chan error, 1)

	_ = webserv.nm.Poller.Start(acceptDesc, func(e netpoll.Event) {

		err := webserv.nm.Pool.ScheduleTimeout(time.Millisecond, func() {
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

		_ = webserv.nm.Poller.Resume(acceptDesc)
	})

	<-webserv.exit

}

func (webserv *WebSocketServer) Stop() {
	webserv.exit <- struct{}{}
}

func (webserv *WebSocketServer) AddHandler(messageType uint16, fun network.RecvHandler) {
	webserv.nm.Group.Hanlder.AddHandler(messageType, fun)
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
