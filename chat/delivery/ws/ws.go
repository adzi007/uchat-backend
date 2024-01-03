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

type chatWebsocket struct {
	Clients              map[*websocket.Conn]bool
	ClientRegisterChanel chan *websocket.Conn
	ClientRemovalChanel  chan *websocket.Conn
	BroadcastChat        chan domain.ChatBubble
	ChatService          domain.ChatService
}

func NewChatRoomHub(cs domain.ChatService) domain.ChatWebsocket {
	return &chatWebsocket{
		Clients:              make(map[*websocket.Conn]bool),
		ClientRegisterChanel: make(chan *websocket.Conn),
		BroadcastChat:        make(chan domain.ChatBubble),
		ClientRemovalChanel:  make(chan *websocket.Conn),
		ChatService:          cs,
	}
}

func (h *chatWebsocket) Run() {
	for {
		select {
		case conn := <-h.ClientRegisterChanel:
			h.Clients[conn] = true

		case conn := <-h.ClientRemovalChanel:
			delete(h.Clients, conn)

		case chat := <-h.BroadcastChat:

			for conn := range h.Clients {
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

func (h *chatWebsocket) Join(client *websocket.Conn) {
	h.ClientRegisterChanel <- client
}

func (h *chatWebsocket) Leave(client *websocket.Conn) {
	h.ClientRemovalChanel <- client
	_ = client.Close()
}

func (h *chatWebsocket) Broadcast(chatBubble domain.ChatBubble) {

	h.BroadcastChat <- chatBubble
}

func (h *chatWebsocket) HandleWsChatRoom() func(*websocket.Conn) {

	return func(c *websocket.Conn) {

		vars := c.Params("chatRoomId")
		room := h.GetRoom(vars)

		if room == nil {
			log.Printf("Room not found")
			c.Close()
			return
		}

		defer room.Leave(c)

		room.Join(c)

		for {

			var chatRequst domain.CreateChatBubbleRequest

			errReadJson := c.ReadJSON(&chatRequst)

			if errReadJson != nil {
				// Handle error
				return
			}

			// INPUT PROSES
			userID, err := primitive.ObjectIDFromHex(chatRequst.UserID)

			if err != nil {

				fmt.Println("gagal", err.Error())
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

			errInserNewChat := h.ChatService.SendChat(chatBubble, chatRequst.ChatRoomId)
			if errInserNewChat != nil {
				fmt.Println("errInserNewChat", errInserNewChat.Error())
				return
			}

			h.Broadcast(chatBubble)

		}

	}
}

func (h *chatWebsocket) GetRoom(roomID string) *chatWebsocket {

	_, err := h.ChatService.GetChatRoomId(roomID)

	if err != nil {
		fmt.Println("gagal 1xxx", err.Error())
		panic(err)
	}

	return h
}
