package delivery

import (
	"adzi-clean-architecture/domain"
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Name  string
	Price string
}

type Hub struct {
	clients              map[*websocket.Conn]bool
	clientRegisterChanel chan *websocket.Conn
	clientRemovalChanel  chan *websocket.Conn
	broadcastChat        chan domain.ChatBubble
	broadcastMessage     chan Message
	chatService          domain.ChatService
}

func NewHub(cs domain.ChatService) *Hub {
	return &Hub{
		clients:              make(map[*websocket.Conn]bool),
		clientRegisterChanel: make(chan *websocket.Conn),
		// broadcastMessage:     make(chan domain.Chat),
		broadcastMessage:    make(chan Message),
		broadcastChat:       make(chan domain.ChatBubble),
		clientRemovalChanel: make(chan *websocket.Conn),
		chatService:         cs,
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

func Chat(h *Hub) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {

		defer func() {
			h.clientRemovalChanel <- conn
			_ = conn.Close()
		}()

		// name := conn.Query("name", "anonymous")
		h.clientRegisterChanel <- conn

		for {

			messageType, _, err := conn.ReadMessage()

			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {

				var chatRequst domain.CreateChatBubbleRequest

				errReadJson := conn.ReadJSON(&chatRequst)

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

				h.broadcastChat <- chatBubble

			}

		}

	}
}

func ChatRoom(h *Hub) func(*websocket.Conn) {

	return func(conn *websocket.Conn) {

		defer func() {
			h.clientRemovalChanel <- conn
			_ = conn.Close()
		}()

		h.clientRegisterChanel <- conn

		fmt.Println("connect to chat room")

		for {

			messageType, pesan, err := conn.ReadMessage()

			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {

				h.broadcastMessage <- Message{
					Price: string(pesan),
				}

			}

			// 	if messageType == websocket.TextMessage {

			// 		var chatRequst domain.CreateChatBubbleRequest

			// 		errReadJson := conn.ReadJSON(&chatRequst)

			// 		if errReadJson != nil {
			// 			// Handle error
			// 			return
			// 		}

			// 		// INPUT PROSES
			// 		userID, err := primitive.ObjectIDFromHex(chatRequst.UserID)

			// 		if err != nil {

			// 			fmt.Println("gagal 1", err.Error())
			// 			panic(err)
			// 		}

			// 		replyId, err := primitive.ObjectIDFromHex(chatRequst.ReplyId)

			// 		chatBubble := domain.ChatBubble{
			// 			ID:        primitive.NewObjectID(),
			// 			Timestamp: time.Now().UTC(),
			// 			UserID:    userID,
			// 			Message:   chatRequst.Message,
			// 			IsDeleted: false,
			// 			ReadedAt:  nil,
			// 		}

			// 		if err == nil {
			// 			chatBubble.ReplyId = replyId
			// 		}

			// 		errInserNewChat := h.chatService.SendChat(chatBubble, chatRequst.ChatRoomId)

			// 		if errInserNewChat != nil {
			// 			fmt.Println("errInserNewChat", errInserNewChat.Error())
			// 			return
			// 		}

			// 		h.broadcastChat <- chatBubble

			// 	}

		}

	}
}
