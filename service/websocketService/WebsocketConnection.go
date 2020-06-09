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
package websocketService

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"net"
	"sort"
	"strconv"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/ice-zzz/netcore/internal/gopool"

	uuid "github.com/satori/go.uuid"
)

type Group struct {
	mu  sync.RWMutex
	seq uint
	us  []*Connection
	ns  map[string]*Connection

	pool *gopool.Pool
	out  chan []byte

	Hanlder *Handler
}

func NewGroup(pool *gopool.Pool) *Group {
	group := &Group{
		pool: pool,
		ns:   make(map[string]*Connection),
		out:  make(chan []byte, 1),
	}
	group.Hanlder = &Handler{}
	group.Hanlder.Init()

	go group.writer()

	return group
}

// Register registers new connection as a User.
func (c *Group) Register(conn net.Conn) *Connection {
	user := &Connection{
		group: c,
		conn:  conn,
		once:  0,
	}

	c.mu.Lock()
	{
		user.id = c.seq
		user.name = c.randName()

		c.us = append(c.us, user)
		c.ns[user.name] = user

		c.seq++
	}
	c.mu.Unlock()

	return user
}

// Remove removes user from chat.
func (c *Group) Remove(user *Connection) {
	c.mu.Lock()
	removed := c.remove(user)
	c.mu.Unlock()

	if !removed {
		return
	}

}

// Rename renames user.
func (c *Group) Rename(user *Connection, name string) (prev string, ok bool) {
	c.mu.Lock()
	{
		if _, has := c.ns[name]; !has {
			ok = true
			prev, user.name = user.name, name
			delete(c.ns, prev)
			c.ns[name] = user
		}
	}
	c.mu.Unlock()

	return prev, ok
}

// writer writes broadcast messages from chat.out channel.
func (c *Group) writer() {
	for bts := range c.out {
		c.mu.RLock()
		us := c.us
		c.mu.RUnlock()

		for _, u := range us {
			u := u // For closure.
			c.pool.Schedule(func() {
				_ = u.writeRaw(bts)
			})
		}
	}
}

// mutex must be held.
func (c *Group) remove(user *Connection) bool {
	if _, has := c.ns[user.name]; !has {
		return false
	}

	delete(c.ns, user.name)

	i := sort.Search(len(c.us), func(i int) bool {
		return c.us[i].id >= user.id
	})
	if i >= len(c.us) {
		// logger.Errorf("Group: 状态不一致")
	}

	without := make([]*Connection, len(c.us)-1)
	copy(without[:i], c.us[:i])
	copy(without[i:], c.us[i+1:])
	c.us = without

	return true
}

func (c *Group) randName() string {
	return uuid.NewV4().String()
}

type Connection struct {
	io   sync.Mutex
	conn io.ReadWriteCloser

	id    uint
	name  string
	group *Group
	once  uint64
}

// Receive reads next message from user's underlying connection.
// It blocks until full message received.
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

	md := u.group.Hanlder.Execute(req)
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
