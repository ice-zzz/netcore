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
package websocketPlugins

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ice-zzz/netcore/service/websocketService"
)

func TestClient_SendMessage(t *testing.T) {
	s := websocketService.New()
	s.AddHandler(0, func(message *websocketService.MessageData) *websocketService.MessageData {

		log.Printf("%s", message.Message)

		return &websocketService.MessageData{
			MessageType: 1,
			Message:     []byte("66666"),
		}

	})
	s.Name = "test"
	s.Ip = "0.0.0.0"
	s.Port = 7777
	go s.Start()

	time.Sleep(time.Second * 3)
	c := NewWebsocketClient("192.168.1.30:7777", "/")
	err := c.Connect()
	if err != nil {
		fmt.Println(err)
	}

	err = c.SendMessage(0, []byte("Hellowrold"))
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Print("")
	}

}
