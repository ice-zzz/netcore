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
package network

import (
	"fmt"
	"log"
	"net"
	"time"

	"git.bitcode.work/ice/netcore/easygo/gopool"
	"git.bitcode.work/ice/netcore/easygo/logs"
	"git.bitcode.work/ice/netcore/easygo/netpoll"
	"github.com/gobwas/ws"
)

type WebSocketOption struct {
	Ip   string         `toml:"ip"`
	Port int            `toml:"port"`
	Name string         `toml:"name"`
	Log  logs.LogOption `toml:"log"`
}

type WebSocketServer struct {
	exit   chan struct{}
	pool   *gopool.Pool
	poller netpoll.Poller
	group  *Group
	logger *logs.Logger
	conf   WebSocketOption
}

type deadliner struct {
	net.Conn
	t time.Duration
}

func CreateWebSocket(opt WebSocketOption) *WebSocketServer {
	pool := gopool.NewPool(128, 1, 1)
	poller, _ := netpoll.New(nil)
	return &WebSocketServer{
		exit:   make(chan struct{}),
		pool:   pool,
		poller: poller,
		group:  NewGroup(pool),
		logger: nil,
		conf:   opt,
	}

}

func (webserv *WebSocketServer) Start() {

	webserv.logger = logs.New(webserv.conf.Log)

	handle := func(conn net.Conn) {

		safeConn := deadliner{conn, time.Millisecond * 100}

		_, err := ws.Upgrade(safeConn)
		if err != nil {
			webserv.logger.Errorf("%s: 升级失败: %v \n", nameConn(conn), err)
			_ = conn.Close()
			return
		}

		user := webserv.group.Register(safeConn)
		desc := netpoll.Must(netpoll.HandleRead(conn))

		webserv.logger.Infof("用户 %s 进入 ip-> %s  \n", user.name, conn.RemoteAddr().String())

		_ = webserv.poller.Start(desc, func(ev netpoll.Event) {
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {

				_ = webserv.poller.Stop(desc)
				webserv.group.Remove(user)
				return
			}

			webserv.pool.Schedule(func() {

				if err := user.Receive(); err != nil {

					_ = webserv.poller.Stop(desc)
					webserv.group.Remove(user)
				}

			})
		})

	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", webserv.conf.Ip, webserv.conf.Port))
	if err != nil {
		log.Fatal(err)
		return
	}

	webserv.logger.Infof("websocket 正在监听端口-> %d \n", webserv.conf.Port)

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		ln, netpoll.EventRead|netpoll.EventOneShot,
	))
	// webserv.EventDispatch.DispatchEvent(COMPLETE, &events.Event{
	// 	EventType: COMPLETE,
	// 	Data:      nil,
	// })

	accept := make(chan error, 1)

	_ = webserv.poller.Start(acceptDesc, func(e netpoll.Event) {

		err := webserv.pool.ScheduleTimeout(time.Millisecond, func() {
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

			webserv.logger.Errorf("连接错误: %v \n", err)

		cooldown:
			delay := 5 * time.Millisecond
			webserv.logger.Errorf("连接错误: %v; %s 秒后重试! \n ", err, delay)
			time.Sleep(delay)
		}

		_ = webserv.poller.Resume(acceptDesc)
	})

	<-webserv.exit

}

func (webserv *WebSocketServer) Stop() {
	webserv.exit <- struct{}{}
}

func (webserv *WebSocketServer) AddHandler(messageType uint16, fun RecvHandler) {
	webserv.group.Hanlder.AddHandler(messageType, fun)
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
