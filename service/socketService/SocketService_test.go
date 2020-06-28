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
	"log"
	"testing"

	"github.com/segmentio/ksuid"
)

func TestSocketServer_Start(t *testing.T) {
	// s := New()
	// s.AddHandler(0, func(message *network.MessageData) *network.MessageData {
	//
	// 	log.Printf("%s", message.Message)
	//
	// 	return &network.MessageData{
	// 		MessageType: 1,
	// 		Message:     []byte("66666"),
	// 	}
	//
	// })
	// s.Name = "test"
	// s.Ip = "0.0.0.0"
	// s.Port = 7777
	// go s.Start()
	//
	// time.Sleep(time.Second*3)
	//
	// server := "192.168.1.30:7777"
	// tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	// 	os.Exit(1)
	// }
	// conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//
	// nm:= network.NewNetManager()
	// nm.Group.Hanlder.AddHandler(1, func(message *network.MessageData) *network.MessageData {
	//
	// 	log.Printf("%s", message.Message)
	//
	// 	return nil
	//
	// })
	// safeConn := deadliner{conn, time.Millisecond * 100}
	// user := nm.Group.Register(safeConn)
	// desc := netpoll.Must(netpoll.HandleRead(conn))
	// _ = nm.Poller.Start(desc, func(ev netpoll.Event) {
	// 	// 断线处理
	// 	if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
	//
	// 		_ = nm.Poller.Stop(desc)
	// 		nm.Group.Remove(user)
	// 		return
	// 	}
	//
	// 	nm.Pool.Schedule(func() {
	//
	// 		if err := user.Receive(); err != nil {
	//
	// 			_ = nm.Poller.Stop(desc)
	// 			nm.Group.Remove(user)
	// 		}
	// 	})
	// })
	//
	// _ =user.Write(&network.MessageData{
	// 	MessageType: 0,
	// 	Message:     []byte("Hello Socket"),
	// })

	for i := 0; i < 20; i++ {
		log.Println(ksuid.New().String())
	}
	// a := ksuid.New().String()
	// b:= ksuid.New().String()
	// c := md5.New()
	// c.Write([]byte(a+b))
	//
	// log.Println(hex.EncodeToString(c.Sum(nil)))

	// for {
	// 	fmt.Print("")
	// }
}
