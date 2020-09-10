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
	"sync"
	"testing"
	"time"
)

func TestTCPClient_Start(t *testing.T) {
	c := &WSClient{
		Mutex:            sync.Mutex{},
		Addr:             "ws://192.168.1.30:9999/",
		ConnNum:          1,
		ConnectInterval:  0,
		PendingWriteNum:  0,
		MaxMsgLen:        64000,
		HandshakeTimeout: 10 * time.Second,
		AutoReconnect:    false,
		NewAgent:         aaaa,
	}

	c.Start()
	time.Sleep(time.Second * 2)

	c.dial().WriteMsg([]byte("hello"))

	for {
		fmt.Print("")
	}
}

func aaaa(*WSConn) Agent {
	return testagent{}
}

type testagent struct {
}

func (t testagent) OnInit() {

}

func (t testagent) React(bytes []byte) {
	log.Printf("%s", bytes)
}

func (t testagent) OnClose() {

}
