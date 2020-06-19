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
	"net"
	"sort"
	"sync"

	"github.com/ice-zzz/netcore/internal/netpoll"
	"github.com/ice-zzz/netcore/utils/gopool"
	"github.com/segmentio/ksuid"
)

type NetManager struct {
	Pool   *gopool.Pool
	Poller netpoll.Poller
	Group  *ConnectionGroup
}

func NewNetManager() *NetManager {
	pool := gopool.NewPool(128, 1, 1)
	poller, _ := netpoll.New(nil)

	return &NetManager{
		Pool:   pool,
		Poller: poller,
		Group:  NewGroup(pool),
	}
}

type ConnectionGroup struct {
	mu  sync.RWMutex
	seq uint
	us  []*Connection
	ns  map[string]*Connection

	pool *gopool.Pool
	out  chan []byte

	Hanlder *Handler
}

func NewGroup(pool *gopool.Pool) *ConnectionGroup {
	group := &ConnectionGroup{
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
func (c *ConnectionGroup) Register(conn net.Conn) *Connection {
	user := &Connection{
		handler: c.Hanlder,
		conn:    conn,
		once:    0,
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
func (c *ConnectionGroup) Remove(user *Connection) {
	c.mu.Lock()
	removed := c.remove(user)
	c.mu.Unlock()

	if !removed {
		return
	}

}

// Rename renames user.
func (c *ConnectionGroup) Rename(user *Connection, name string) (prev string, ok bool) {
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
func (c *ConnectionGroup) writer() {
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
func (c *ConnectionGroup) remove(user *Connection) bool {
	if _, has := c.ns[user.name]; !has {
		return false
	}

	delete(c.ns, user.name)

	i := sort.Search(len(c.us), func(i int) bool {
		return c.us[i].id >= user.id
	})
	if i >= len(c.us) {
		// logger.Errorf("ConnectionGroup: 状态不一致")
	}

	without := make([]*Connection, len(c.us)-1)
	copy(without[:i], c.us[:i])
	copy(without[i:], c.us[i+1:])
	c.us = without

	return true
}

func (c *ConnectionGroup) randName() string {
	return ksuid.New().String()
}

type Handler struct {
	funList map[uint16]RecvHandler
}

type RecvHandler func(message *MessageData) *MessageData

func (h *Handler) Init() {
	h.funList = make(map[uint16]RecvHandler)
}

func (h *Handler) AddHandler(messageType uint16, fun RecvHandler) {
	h.funList[messageType] = fun
}

type MessageData struct {
	MessageType uint16
	Message     []byte
}

func (h *Handler) Execute(data *MessageData) *MessageData {
	if v, ok := h.funList[data.MessageType]; ok {
		return v(data)
	}
	return &MessageData{
		MessageType: 500,
		Message:     []byte("非法消息类型"),
	}
}
