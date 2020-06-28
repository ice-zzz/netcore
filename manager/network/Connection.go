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
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
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
	if md == nil {
		return nil
	}
	err = u.write(md)
	if err != nil {
		_ = u.conn.Close()
		return err
	}
	return nil
}

func (u *Connection) Write(md *MessageData) error {
	err := u.write(md)
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

	var r io.Reader
	var err error

	if strings.Split(u.name, "_")[0] == "ws" {
		var h ws.Header
		h, r, err = wsutil.NextReader(u.conn, ws.StateServerSide)
		if err != nil {
			return nil, err
		}
		if h.OpCode.IsControl() {
			return nil, wsutil.ControlFrameHandler(u.conn, ws.StateServerSide)(h, r)
		}
	} else {
		r = u.conn

	}

	// var bs []byte
	// bs, err = read(r)

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

// 包格式为 0xFF|0xFF|token(s 27bit)(r 32bit)|包总数|当前数|len(高)|len(低)|MsgType|Message(zlib)|0xFF|0xFE
// 0xFF|0xFF 起始标识符
// token(32bit) 为随机字符,客户端需要返回 token(27)+uid(27) 并且md5(32)
// 其中len为data的长度，实际长度为len(高)*256+len(低)
//
func read(reader io.Reader) ([]byte, error) {

	recvCache := make(map[uint64][][]byte)
	once := uint64(0)

	// 状态机状态
	state := 0x00
	// 数据包长度
	length := uint16(0)
	msgid := uint64(0)
	packindex := uint16(0)
	// msgtype := uint16(0)

	var (
		recvBuffer, shaBuffer, tokenBuffer []byte
	)
	// 游标
	cursor := uint16(0)

	defer func() {
		recover()
	}()

	for {
		reader := bufio.NewReader(reader)
		for {
			recvByte, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					// logger.Infof("%s 用户退出", sc.conn.RemoteAddr())
				}
				return nil, err
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
				if cursor == 32 {
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
				// msgtype = uint16(recvByte)
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
						// sc.messageInChannel <- &network.Work{
						// 	OutChannel: sc.messageOutChannel,
						// 	Message: &network.MessageData{
						// 		MessageType: msgtype,
						// 		Message:     recvBuffer,
						// 	},
						// }
					}
				}
				// 状态机归位,接收下一个包

				cursor = 0
				state = 0x00
				break
			}

		}
	}

}

func zip(inData []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, _ = w.Write(inData)
	_ = w.Close()
	return in.Bytes()
}

func checkFull(packs [][]byte) bool {
	for _, v := range packs {
		if v == nil {
			return false
		}
	}
	return true

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

// func zip(inData []byte) []byte {
// 	var in bytes.Buffer
// 	w := zlib.NewWriter(&in)
// 	_, _ = w.Write(inData)
// 	_ = w.Close()
// 	return in.Bytes()
// }

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
