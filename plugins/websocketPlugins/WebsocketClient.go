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
	"net/url"
	"strconv"

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
