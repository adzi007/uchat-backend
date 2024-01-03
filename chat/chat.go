package chat

import (
	"adzi-clean-architecture/chat/delivery/http"
	delivery "adzi-clean-architecture/chat/delivery/ws"
	"adzi-clean-architecture/chat/repository"
	"adzi-clean-architecture/chat/service"
	"adzi-clean-architecture/config"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func InitChat(app *fiber.App, dbConfig config.MongoDbInterface) {

	// Repository
	chatRoom := dbConfig.GetCollection("chatRoom")
	// userCollection := dbConfig.GetCollection("chatRoom")

	chatRepo := repository.NewChatRepo(chatRoom)

	// Service
	chatService := service.NewChatService(chatRepo)

	// hub := delivery.NewHub(chatService)
	hub := delivery.NewChatRoomHub(chatService)
	go hub.Run()

	app.Use("/ws", delivery.AllowUpgrade)
	// app.Use("/ws/chat/:chatRoomId", websocket.New(delivery.Chat(hub)))
	app.Use("/ws/chat/:chatRoomId", websocket.New(hub.HandleWsChatRoom()))

	http.NewRouteUser(app, chatService, hub)

}
