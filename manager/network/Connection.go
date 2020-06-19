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
package network

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

const (
	DataSliceLength = 1400
)

type Connection struct {
	io   sync.Mutex
	conn io.ReadWriteCloser

	id      uint
	name    string
	handler *Handler
	once    uint64
}

func (u *Connection) GetName() string {
	return u.name
}

// 接收从用户的连接读取下一条消息。
// 它会阻塞直到收到完整的消息。
func (u *Connection) Receive() error {
	req, err := u.readRequest()
	if err != nil {
		_ = u.conn.Close()
		return err
	}
	if req == nil {
		// Handled some control message.
		return nil
	}

	md := u.handler.Execute(req)
	err = u.write(md)
	if err != nil {
		_ = u.conn.Close()
		return err
	}
	return nil
}

// 	包格式为 once(16bit 0-f)|MsgType |Message(RSA)

func (u *Connection) readRequest() (*MessageData, error) {
	u.io.Lock()
	defer u.io.Unlock()

	h, r, err := wsutil.NextReader(u.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(u.conn, ws.StateServerSide)(h, r)
	}

	req := &MessageData{}
	message, err := ioutil.ReadAll(r)
	once, err := strconv.ParseUint(string(message[:16]), 16, 32)
	msgtype := uint16(message[16:17][0])
	if once > u.once {
		u.once = once
	} else {
		return nil, nil
	}

	//
	// if err != nil {
	// 	return nil, err
	// }
	// req.Message, err = crypto.RSA.PriKeyDECRYPT(message[3:])
	// if err != nil {
	// 	logger.Errorf("RSA Decode Error")
	// 	return nil, err
	// }
	req.MessageType = msgtype
	req.Message = message[17:]
	return req, nil

}

func (u *Connection) write(md *MessageData) error {
	w := wsutil.NewWriter(u.conn, ws.StateServerSide, ws.OpBinary)

	u.io.Lock()
	defer u.io.Unlock()

	writeByte := make([]byte, 0)
	writeByte = append(writeByte, byte(u.once))
	writeByte = append(writeByte, byte(md.MessageType))
	message := md.Message
	// message, err := crypto.RSA.PriKeyENCTYPT(md.Message)
	// if err != nil {
	// 	logger.Errorf("RSA Encode Error")
	// 	return err
	// }
	writeByte = append(writeByte, []byte(message)...)

	_, _ = w.Write(writeByte)

	return w.Flush()
}

func (u *Connection) writeRaw(p []byte) error {
	u.io.Lock()
	defer u.io.Unlock()

	_, err := u.conn.Write(p)

	return err
}

func zip(inData []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, _ = w.Write(inData)
	_ = w.Close()
	return in.Bytes()
}

func (u *Connection) EnCoder(message []byte, messageType uint16) [][]byte {

	start := 0
	end := 0
	zippedData := message
	// zippedData := zip(message)
	zippedDataLength := len(zippedData)
	zippedDataSliceNum := int(math.Ceil(float64(zippedDataLength) / float64(DataSliceLength)))
	zippedDataSliceNum = int(math.Max(float64(zippedDataSliceNum), 1))

	sendBytes := make([][]byte, zippedDataSliceNum)
	for i := 0; i < zippedDataSliceNum; i++ {

		if i*DataSliceLength < zippedDataLength {
			if (i+1)*DataSliceLength-1 > zippedDataLength {
				start = i * DataSliceLength
				end = zippedDataLength
			} else {
				start = i * DataSliceLength
				end = (i + 1) * DataSliceLength
			}
		}
		zippedSendData := zippedData[start:end]
		sendBytes[i] = append(sendBytes[i], byte(0xFF))
		sendBytes[i] = append(sendBytes[i], byte(0xFF))
		sendBytes[i] = append(sendBytes[i], byte(u.once>>8))
		sendBytes[i] = append(sendBytes[i], byte(u.once&0xFF))
		// sendBytes[i] = BytesCombine(sendBytes[i], []byte(sc.connId))// TODO 有状况
		sendBytes[i] = append(sendBytes[i], byte(zippedDataSliceNum))
		sendBytes[i] = append(sendBytes[i], byte(i))

		currentLength := len(zippedSendData)
		sendBytes[i] = append(sendBytes[i], byte(uint16(currentLength)>>8))
		sendBytes[i] = append(sendBytes[i], byte(uint16(currentLength)&0xFF))

		sha := sha256.New()
		sha.Write(zippedSendData)
		code := sha.Sum([]byte(nil))
		sendBytes[i] = append(sendBytes[i], code[:16]...)
		sendBytes[i] = append(sendBytes[i], byte(messageType))
		sendBytes[i] = append(sendBytes[i], zippedSendData...)
		sendBytes[i] = append(sendBytes[i], byte(0xFF))
		sendBytes[i] = append(sendBytes[i], byte(0xFE))
	}

	return sendBytes
}
