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
	"os"

	"github.com/panjf2000/gnet"
	"github.com/sirupsen/logrus"
)

type HttpService struct {
	*gnet.EventServer
	Addr        string
	exitChannel chan os.Signal
	logger      *logrus.Logger
}

func (hs *HttpService) Start() {
}

func (hs *HttpService) Stop() {
}

func (hs *HttpService) OnInitComplete(svr gnet.Server) (action gnet.Action) {
	return
}

func (hs *HttpService) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	return
}

func (hs *HttpService) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	return
}

func (hs *HttpService) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	return
}
