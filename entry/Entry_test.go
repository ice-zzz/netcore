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
package entry

import (
	"log"
	"testing"
	"time"

	"github.com/ice-zzz/netcore/service/websocketService"
)

func TestEntry_Start(t *testing.T) {
	e := Create()

	wss := websocketService.New()
	defer wss.Stop()
	wss.AddHandler(0, func(message *websocketService.MessageData) *websocketService.MessageData {

		log.Printf("%s", message.Message)

		return &websocketService.MessageData{
			MessageType: 1,
			Message:     []byte("66666"),
		}

	})
	wss.Name = "test"
	wss.Ip = "0.0.0.0"
	wss.Port = 7777
	e.AddService(wss)
	go e.Start()
	time.Sleep(time.Second * 3)

	log.Println(e.services[wss.GetServiceName()].IsRunning())

}
