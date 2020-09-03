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
	"io"
	"log"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type DataHandlerParam struct {
	Request []byte

	OpCode ws.OpCode

	Writer io.Writer

	WSConn *GnetUpgraderConn // 升级后的连接

	Server *WebSocket
}

type DataHandler func(param *DataHandlerParam)

// 简单的echo server
func EchoDataHandler(param *DataHandlerParam) {

	log.Printf("server 接收到数据, opcode:%x, %s\n", param.OpCode, string(param.Request))

	response := fmt.Sprintf("response is :%s, 当前时间:%s\n", string(param.Request), time.Now().Format("2006-01-02 15:04:05"))

	// param.Writer.Write([]byte(response))

	// ws.WriteFrame(param.Writer, ws.NewTextFrame([]byte(response)))

	_ = wsutil.WriteServerMessage(param.Writer, param.OpCode, []byte(response))

	return
}
