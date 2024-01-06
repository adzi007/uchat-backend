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

type roomChanel struct {
	Client *websocket.Conn
	RoomId string
}

type broadcastChanel struct {
	RoomId      string
	ChatMessage domain.ChatBubble
}

type chatWebsocket struct {
	Clients              map[*websocket.Conn]bool
	ClientRegisterChanel chan *websocket.Conn
	ClientRemovalChanel  chan *websocket.Conn
	BroadcastChat        chan broadcastChanel
	ChatService          domain.ChatService
	ChatRoomConnections  []roomChanel
}

func NewChatRoomHub(cs domain.ChatService) domain.ChatWebsocket {
	return &chatWebsocket{
		Clients:              make(map[*websocket.Conn]bool),
		ClientRegisterChanel: make(chan *websocket.Conn),
		BroadcastChat:        make(chan broadcastChanel),
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

			for i := range h.ChatRoomConnections {

				if h.ChatRoomConnections[i].RoomId == chat.RoomId {
					_ = h.ChatRoomConnections[i].Client.WriteJSON(chat.ChatMessage)
				}

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

func (h *chatWebsocket) Join(client *websocket.Conn, roomId string) {

	h.ClientRegisterChanel <- client

	h.ChatRoomConnections = append(h.ChatRoomConnections, roomChanel{
		Client: client,
		RoomId: roomId,
	})

}

func (h *chatWebsocket) Leave(client *websocket.Conn) {
	h.ClientRemovalChanel <- client
	_ = client.Close()
}

func (h *chatWebsocket) Broadcast(chatBubble domain.ChatBubble, roomId string) {

	h.BroadcastChat <- broadcastChanel{
		RoomId:      roomId,
		ChatMessage: chatBubble,
	}
}

func (h *chatWebsocket) HandleWsChatRoom() func(*websocket.Conn) {

	return func(c *websocket.Conn) {

		chatRoomId := c.Params("chatRoomId")
		room := h.GetRoom(chatRoomId)

		if room == nil {
			log.Printf("Room not found")
			c.Close()
			return
		}

		defer room.Leave(c)

		room.Join(c, chatRoomId)

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

			errInserNewChat := h.ChatService.SendChat(chatBubble, chatRoomId)
			if errInserNewChat != nil {
				fmt.Println("errInserNewChat", errInserNewChat.Error())
				return
			}

			h.Broadcast(chatBubble, chatRoomId)

		}

	}
}

func (h *chatWebsocket) GetRoom(roomID string) *chatWebsocket {

	_, err := h.ChatService.GetChatRoomId(roomID)

	if err != nil {
		fmt.Println("gagal", err.Error())
		panic(err)
	}

	return h
}
