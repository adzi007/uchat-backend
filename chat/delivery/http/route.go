package http

import (
	"adzi-clean-architecture/config"
	"adzi-clean-architecture/domain"
	"adzi-clean-architecture/pkg/utils"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	Cs     domain.ChatService
	ChatWs domain.ChatWebsocket
}

func NewRouteUser(r *fiber.App, us domain.ChatService, chatWs domain.ChatWebsocket) {

	handler := &ChatHandler{
		Cs:     us,
		ChatWs: chatWs,
	}

	r.Static("/public", config.ProjectRootPath+"/public")

	r.Post("/chat/new", handler.CreateNewChat)
	r.Post("/chat/send", handler.SendChat)
	r.Put("/chat/setReaded", handler.SetReaded)

}

func (uh *ChatHandler) CreateNewChat(ctx *fiber.Ctx) error {

	chatRequst := new(domain.ChatCreateRequest)

	if err := ctx.BodyParser(chatRequst); err != nil {
		return err
	}

	// VALIDASI INPUT

	var validate = validator.New()

	errValidate := validate.Struct(chatRequst)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	// HANDLE FILE
	attachment := handleMediaUpload(ctx)

	// INPUT PROSES

	UserCreatId, err := primitive.ObjectIDFromHex(chatRequst.UserCreate)

	if err != nil {
		panic(err)
	}

	UserReceiverId, err := primitive.ObjectIDFromHex(chatRequst.UserReceiver)

	if err != nil {
		panic(err)
	}

	newChat := domain.Chat{
		ID:          primitive.NewObjectID(),
		TimeCreated: time.Now().UTC(),
		Members:     []primitive.ObjectID{UserCreatId, UserReceiverId},
		Messages: []domain.ChatBubble{
			{
				ID:         primitive.NewObjectID(),
				Timestamp:  time.Now().UTC(),
				UserID:     UserCreatId,
				Message:    chatRequst.Message,
				Attachment: attachment,
				IsDeleted:  false,
			},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	errInserNewChat := uh.Cs.CreateNewChat(newChat)

	if errInserNewChat != nil {

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "failed inser new chat",
			"error": errInserNewChat.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "success create new chat",
		"data":  newChat,
	})
}

func handleMediaUpload(ctx *fiber.Ctx) domain.ChatBubbleMedia {

	var attachment domain.ChatBubbleMedia

	if form, err := ctx.MultipartForm(); err == nil {

		files := form.File["attachment"]

		// Loop through files:
		for _, file := range files {

			fileType := utils.CheckFileType(file.Filename)

			ext := utils.GetFileExtension(file.Filename)

			// filename = file.Filename
			filename := utils.GenerateUniqueFileName(ext)

			errSaveFile := ctx.SaveFile(file, fmt.Sprintf("./public/attachment/%s", filename))

			if errSaveFile != nil {
				log.Println("fail to store file into public/attachment")
			}

			if fileType == "Image" {

				attachment.Images = append(attachment.Images, filename)

			} else if fileType == "Video" {

				attachment.Video = filename
			} else {

				attachment.Document.FileName = filename
				attachment.Document.Size = int32(file.Size)

			}

		}

	}

	return attachment
}

func (uh *ChatHandler) SendChat(ctx *fiber.Ctx) error {

	chatRequst := new(domain.CreateChatBubbleRequest)

	if err := ctx.BodyParser(chatRequst); err != nil {
		return err
	}

	// VALIDASI INPUT

	var validate = validator.New()

	errValidate := validate.Struct(chatRequst)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed validate",
			"error":   errValidate.Error(),
		})
	}

	// END HANDLE FILE ---------------------------------------------------

	attachment := handleMediaUpload(ctx)

	// INPUT PROSES
	userID, err := primitive.ObjectIDFromHex(chatRequst.UserID)

	if err != nil {
		panic(err)
	}

	replyId, err := primitive.ObjectIDFromHex(chatRequst.ReplyId)

	if err != nil {
		fmt.Println(err.Error())

	}

	chatBubble := domain.ChatBubble{
		ID:         primitive.NewObjectID(),
		Timestamp:  time.Now().UTC(),
		UserID:     userID,
		ReplyId:    replyId,
		Message:    chatRequst.Message,
		Attachment: attachment,
		IsDeleted:  false,
	}

	errInserNewChat := uh.Cs.SendChat(chatBubble, chatRequst.ChatRoomId)

	if errInserNewChat != nil {

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "failed inser new chat",
			"error": errInserNewChat.Error(),
		})
	}

	uh.ChatWs.Broadcast(chatBubble)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "success create new chat",
		"data":  chatBubble,
	})
}

func (uh *ChatHandler) SetReaded(ctx *fiber.Ctx) error {

	readedRequest := new(domain.SetReadedRequest)

	if err := ctx.BodyParser(readedRequest); err != nil {
		return err
	}

	// VALIDASI INPUT

	var validate = validator.New()

	errValidate := validate.Struct(readedRequest)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed validate",
			"error":   errValidate.Error(),
		})
	}

	// UPDATE PROSES

	err := uh.Cs.SetReadedChat(readedRequest.ChatRoomId, readedRequest.ChatBubbleId)

	if err != nil {

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "failed set readed chat by userId " + readedRequest.UserID,
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "success set readed status by userId " + readedRequest.UserID,
	})
}
