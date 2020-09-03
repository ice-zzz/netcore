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
package http

import (
	"github.com/panjf2000/gnet"
)

type HttpCode struct {
}

func (h *HttpCode) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	panic("implement me")
}

func (h *HttpCode) Decode(c gnet.Conn) (data []byte, err error) {
	data = c.Read()
	switch svr := c.Context().(type) {
	case nil: //   首次进入
		c.SetContext(&Context{
			req:  Request{Header: map[string]string{}},
			conn: nil,
			Head: nil,
			Data: nil,
		})
		return h.Decode(c)
	case *Context:
		shiftN, data, err := svr.req.Parsereq(data)
		c.ShiftN(shiftN)
		return data, err
	}

	return
}
