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
package websocket

import (
	"fmt"
	"testing"

	"github.com/ice-zzz/netcore/services/http"
	"github.com/panjf2000/gnet"
)

func TestNewEchoServer(t *testing.T) {
	port := 9999

	tcpServer := NewEchoServer(fmt.Sprintf("0.0.0.0:%d", port))

	gnet.Serve(tcpServer,
		fmt.Sprintf("tcp://:%d", port),
		gnet.WithMulticore(true),
		gnet.WithCodec(&http.HttpCode{}),
	)

}

func TestNewServer(t *testing.T) {

	// conn, err := net.Dial("tcp", "192.168.1.30:9999")
	// conn.Write([]byte("AAAAA"))
	// fmt.Println(err)
	// for {
	//
	// }
}
