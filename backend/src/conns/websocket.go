package conns

import (
	"context"
	"log/slog"
	"sync"
	"zzc/fall-chat/src/keys/ctxkey"
	"zzc/fall-chat/src/utils"

	"github.com/gofiber/contrib/v3/websocket"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte // 发送通道，防止并发写
	logger *slog.Logger
	ctx    context.Context
}

func NewClient(ctx context.Context, id string, conn *websocket.Conn) *Client {
	return &Client{
		ctx:  ctx,
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 1024),
	}
}

func (c *Client) GetContext() context.Context {
	return c.ctx
}

func (c *Client) GetLogger() *slog.Logger {
	if c.logger == nil {
		c.logger = ctxkey.GetLogger(c.ctx)
	}
	return c.logger
}

func (c *Client) LogError(err error, fields ...any) {
	fields = append(fields, "error", err.Error())
	c.GetLogger().Error("websocket", fields...)
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		// 并发安全的写
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}

func (c *Client) SendMessage(msg []byte) {
	select {
	case c.Send <- msg:
	default:
		// 通道满，客户端消费慢或断网了
		utils.LogWarn("客户端 %s 发送缓冲区已满，丢弃消息\n", c.ID)
		// 踢掉它
		c.Conn.Close()
	}
}

type Hub struct {
	clients map[string]*Client
	groups  map[string]map[string]bool
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		groups:  make(map[string]map[string]bool),
	}
}

func (h *Hub) AddClient(ctx context.Context, id string, conn *websocket.Conn) *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	client := NewClient(ctx, id, conn)
	h.clients[id] = client
	// 启动写通道
	go client.WritePump()

	return client
}

func (h *Hub) AddGroups(id string, groupIds []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, groupId := range groupIds {
		if h.groups[groupId] == nil {
			h.groups[groupId] = make(map[string]bool)
		}
		h.groups[groupId][id] = true
	}
}

func (h *Hub) AddGroup(id string, groupID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.groups[groupID] == nil {
		h.groups[groupID] = make(map[string]bool)
	}
	h.groups[groupID][id] = true
}

func (h *Hub) RemoveClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, id)
	h.removeGroup(id)
}

func (h *Hub) removeGroup(id string) {
	for gID, group := range h.groups {
		delete(group, id)

		if len(group) == 0 {
			delete(h.groups, gID)
		}
	}
}

func (h *Hub) GetClient(id string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, exists := h.clients[id]
	return client, exists
}

func (h *Hub) GetGroup(id string) (map[string]bool, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	group, exists := h.groups[id]
	return group, exists
}
