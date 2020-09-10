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
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ice-zzz/netcore/gate"
)

func main() {
	g := gate.Gate{
		MaxConnNum:     2000,
		MaxMsgLen:      64000,
		Processor:      nil,
		WSAddr:         "0.0.0.0:9999",
		HTTPTimeout:    15 * time.Second,
		TCPAddr:        "",
		LenMsgLen:      2,
		LittleEndian:   false,
		LenMsgLenInMsg: false,
		ChanStop:       false,
	}
	gclose := make(chan bool)
	g.Run(gclose)

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	gclose <- true
	os.Exit(1)

}
