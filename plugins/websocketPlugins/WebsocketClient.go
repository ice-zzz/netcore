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
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/ice-zzz/netcore/manager/network"
	"golang.org/x/net/websocket"
)

type Client struct {
	Host string
	Path string
	ws   *websocket.Conn
}

func NewWebsocketClient(host, path string) *Client {
	return &Client{
		Host: host,
		Path: path,
	}
}

func (c *Client) Connect() error {
	u := url.URL{Scheme: "ws", Host: c.Host, Path: c.Path}

	ws, err := websocket.Dial(u.String(), "", "http://"+c.Host+"/")
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	c.ws = ws

	go func() {

		recvCache := make(map[uint64][][]byte)
		once := uint64(0)

		// 状态机状态
		state := 0x00
		// 数据包长度
		length := uint16(0)
		// msgtype := uint16(0)
		msgid := uint64(0)
		packindex := uint16(0)
		msgtype := uint16(0)

		var (
			recvBuffer, shaBuffer, tokenBuffer []byte
		)
		// 游标
		cursor := uint16(0)

		defer func() {
			recover()
		}()

		for {
			reader := bufio.NewReader(c.ws)
			for {
				recvByte, err := reader.ReadByte()
				if err != nil {
					if err == io.EOF {
						fmt.Printf("%s 用户退出", c.ws.RemoteAddr())
					}
					return
				}

				switch state {
				case 0x00:
					if recvByte == 0xFF {
						state = 0x01
						// 初始化状态机
						recvBuffer = nil
						length = 0
					} else {
						state = 0x00
					}
					break
				case 0x01:
					if recvByte == 0xFF {
						tokenBuffer = make([]byte, 40)
						state = 0x02
					} else {
						state = 0x00
					}
					break

				case 0x02:
					once = 0
					once += uint64(recvByte) * 256
					state = 0x03
					break
				case 0x03:
					once += uint64(recvByte)
					state = 0x04
					break

				case 0x04:
					tokenBuffer[cursor] = recvByte
					cursor++
					if cursor == 40 {
						// if sc.connId == fmt.Sprintf("%x", tokenBuffer) {
						state = 0x05
						// } else {
						// 	state = 0x00
						// }
						cursor = 0
					}
					break
				case 0x05:
					t := uint16(recvByte)
					if _, ok := recvCache[msgid]; !ok {
						recvCache[msgid] = make([][]byte, t)
					}
					state = 0x06
					break
				case 0x06:
					packindex = uint16(recvByte)
					state = 0x07
					break
				case 0x07:
					length += uint16(recvByte) * 256
					state = 0x08
					break
				case 0x08:
					length += uint16(recvByte)
					// 一次申请缓存，初始化游标，准备读数据
					recvBuffer = make([]byte, length)
					shaBuffer = make([]byte, 16)
					cursor = 0
					state = 0x09
					break

				case 0x09:
					shaBuffer[cursor] = recvByte
					cursor++
					if cursor == 16 {
						state = 0x0A
						cursor = 0
					}
					break
				case 0x0A:
					msgtype = uint16(recvByte)
					state = 0x0B
					break
				case 0x0B:

					// 不断地在这个状态下读数据，直到满足长度为止
					recvBuffer[cursor] = recvByte
					cursor++
					if cursor >= length {
						recvCache[msgid][packindex] = recvBuffer
						state = 0x0C
					}
				case 0x0C:
					if recvByte == 0xFF {
						state = 0x0D
					} else {
						state = 0x00
					}
					break
				case 0x0D:
					if recvByte == 0xFE {
						sha := sha256.New()
						sha.Write(recvBuffer)
						code := sha.Sum(nil)
						originalSHA := hex.EncodeToString(shaBuffer)
						currentSHA := hex.EncodeToString(code[:16])
						if originalSHA == currentSHA {
							if !checkFull(recvCache[msgid]) {
								state = 0x00
								break
							}
							recvBuffer = make([]byte, 0)
							for _, v := range recvCache[msgid] {
								recvBuffer = append(recvBuffer, v...)
							}

							// zip
							// zlibData := bytes.NewReader(recvBuffer)
							// zlibReader, err := zlib.NewReader(zlibData)
							// if err != nil {
							// 	break
							// }
							// recvBuffer, err = ioutil.ReadAll(zlibReader)
							// if err != nil {
							// 	log.Printf("err: %s", err.Error())
							// 	break
							// }

							a := &network.MessageData{
								MessageType: msgtype,
								Message:     recvBuffer,
							}
							fmt.Println(a)

						}
					}
					// 状态机归位,接收下一个包

					cursor = 0
					state = 0x00
					break
				}

			}
		}

	}()

	return nil
}

func (c *Client) SendMessage(mt uint16, body []byte) error {

	hexstr := strconv.FormatUint(1, 16)
	if l := 16 - len(hexstr); l != 0 {
		for i := 0; i < l; i++ {
			hexstr = "0" + hexstr
		}
	}
	data := make([]byte, 0)
	data = append(data, []byte(hexstr)...)
	data = append(data, byte(mt))
	data = append(data, body...)
	_, err := c.ws.Write(data)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	return nil
}

func (c *Client) Close() error {
	return c.ws.Close()
}

func checkFull(packs [][]byte) bool {
	for _, v := range packs {
		if v == nil {
			return false
		}
	}
	return true

}
