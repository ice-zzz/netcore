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
package codec

import (
	"github.com/panjf2000/gnet"
)

type WSCode struct {
}

func (code *WSCode) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

func (code *WSCode) Decode(c gnet.Conn) (data []byte, err error) {

	switch svr := c.Context().(type) {
	case *WSconn:
		data = c.Read()
		msgtype, data, err := svr.ReadMessage(data)
		c.ShiftN(svr.ReadLength)
		if msgtype == BinaryMessage {
			return data, nil
		}
		return nil, err
	case *Httpserver:
		data = c.Read()
		shift, data, err := svr.Request.Parsereq(data)
		c.ShiftN(shift)
		return data, err
	case nil:

		http := Httppool.Get().(*Httpserver)
		http.c = c
		c.SetContext(http)
		return code.Decode(c)
	}
	return
}
