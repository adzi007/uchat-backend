package delivery

import (
	"adzi-clean-architecture/domain"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRoomHub struct {
	clients              map[*websocket.Conn]bool
	clientRegisterChanel chan *websocket.Conn
	clientRemovalChanel  chan *websocket.Conn
	broadcastChat        chan domain.ChatBubble
	chatService          domain.ChatService
}

func NewChatRoomHub(cs domain.ChatService) *ChatRoomHub {
	return &ChatRoomHub{
		clients:              make(map[*websocket.Conn]bool),
		clientRegisterChanel: make(chan *websocket.Conn),
		// broadcastMessage:     make(chan Message),
		broadcastChat:       make(chan domain.ChatBubble),
		clientRemovalChanel: make(chan *websocket.Conn),
		chatService:         cs,
	}
}

func (h *ChatRoomHub) Run() {
	for {
		select {
		case conn := <-h.clientRegisterChanel:
			h.clients[conn] = true

		case conn := <-h.clientRemovalChanel:
			delete(h.clients, conn)

		case chat := <-h.broadcastChat:
			for conn := range h.clients {
				_ = conn.WriteJSON(chat)
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

func (h *ChatRoomHub) Join(client *websocket.Conn) {
	h.clientRegisterChanel <- client
}

func (h *ChatRoomHub) Leave(client *websocket.Conn) {
	h.clientRemovalChanel <- client
	_ = client.Close()
}

func (h *ChatRoomHub) Broadcast(chatBubble domain.ChatBubble) {

	h.broadcastChat <- chatBubble
}

func HandleWsChatRoom(h *ChatRoomHub) func(*websocket.Conn) {

	return func(c *websocket.Conn) {

		vars := c.Params("chatRoomId")
		room := h.GetRoom(vars)

		if room == nil {
			log.Printf("Room not found")
			c.Close()
			return
		}

		room.Join(c)

		defer room.Leave(c)

		for {

			messageType, _, err := c.ReadMessage()

			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {

				var chatRequst domain.CreateChatBubbleRequest

				errReadJson := c.ReadJSON(&chatRequst)

				if errReadJson != nil {
					// Handle error
					return
				}

				// INPUT PROSES
				userID, err := primitive.ObjectIDFromHex(chatRequst.UserID)

				if err != nil {

					fmt.Println("gagal 1", err.Error())
					panic(err)
				}

				replyId, err := primitive.ObjectIDFromHex(chatRequst.ReplyId)

				chatBubble := domain.ChatBubble{
					ID:        primitive.NewObjectID(),
					Timestamp: time.Now().UTC(),
					UserID:    userID,
					Message:   chatRequst.Message,
					IsDeleted: false,
					ReadedAt:  nil,
				}

				if err == nil {
					chatBubble.ReplyId = replyId
				}

				errInserNewChat := h.chatService.SendChat(chatBubble, chatRequst.ChatRoomId)

				if errInserNewChat != nil {
					fmt.Println("errInserNewChat", errInserNewChat.Error())
					return
				}

				// h.broadcastChat <- chatBubble

				// room.Broadcast(chatBubble)
				h.Broadcast(chatBubble)

			}

		}

	}
}

func (h *ChatRoomHub) GetRoom(roomID string) *ChatRoomHub {

	_, err := h.chatService.GetChatRoomId(roomID)

	if err != nil {
		fmt.Println("gagal 1xxx", err.Error())
		panic(err)
	}

	return h
}
