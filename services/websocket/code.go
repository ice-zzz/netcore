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
	"github.com/panjf2000/gnet"
)

const (
	HttpHead = "GET / HTTP/1.1"
)

type BigCode struct {
}

func (b BigCode) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	panic("implement me")
}

func (b BigCode) Decode(c gnet.Conn) (data []byte, err error) {
	data = c.Read()

	switch c.Context().(type) {
	case nil: //   首次进入

		return
	}

	return
}
