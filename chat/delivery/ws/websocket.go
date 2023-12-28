package delivery

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Name  string
	Price string
}

type Hub struct {
	clients              map[*websocket.Conn]bool
	clientRegisterChanel chan *websocket.Conn
	clientRemovalChanel  chan *websocket.Conn
	// broadcastMessage     chan domain.Chat
	broadcastMessage chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:              make(map[*websocket.Conn]bool),
		clientRegisterChanel: make(chan *websocket.Conn),
		// broadcastMessage:     make(chan domain.Chat),
		broadcastMessage:    make(chan Message),
		clientRemovalChanel: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.clientRegisterChanel:
			h.clients[conn] = true

		case conn := <-h.clientRemovalChanel:
			delete(h.clients, conn)

		case msg := <-h.broadcastMessage:
			for conn := range h.clients {
				_ = conn.WriteJSON(msg)
			}
		}
	}
}

func AllowUpgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}

	return fiber.ErrUpgradeRequired
}

func Chat(h *Hub) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {

		defer func() {
			h.clientRemovalChanel <- conn
			_ = conn.Close()
		}()

		name := conn.Query("name", "anonymous")
		h.clientRegisterChanel <- conn

		for {
			messageType, pesannya, err := conn.ReadMessage()

			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {
				h.broadcastMessage <- Message{
					Name:  name,
					Price: string(pesannya),
				}
			}

		}

	}
}
