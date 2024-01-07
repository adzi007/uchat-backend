package service

import (
	"adzi-clean-architecture/domain"
	"time"
)

type chatService struct {
	chatRepo domain.ChatRepository
	userRepo domain.UserRepository
}

func NewChatService(cr domain.ChatRepository, ur domain.UserRepository) domain.ChatService {
	return &chatService{
		chatRepo: cr,
		userRepo: ur,
	}
}

func (cs chatService) CreateNewChat(newChat domain.NewChatRequest) (domain.Chat, error) {

	userCreatedObj, errGetUserCreated := cs.userRepo.GetUserById(newChat.Members[0])

	if errGetUserCreated != nil {
		panic(errGetUserCreated)
	}

	userReceive, errGetUserReceiveer := cs.userRepo.GetUserById(newChat.Members[1])

	if errGetUserReceiveer != nil {
		panic(errGetUserReceiveer)
	}

	newChatInsert := domain.Chat{
		ID:          newChat.ID,
		TimeCreated: newChat.TimeCreated,
		Members: map[string]*domain.UserResponseData{
			newChat.Members[0]: {
				Id:    userCreatedObj.Id,
				Nama:  userCreatedObj.Nama,
				Phone: userCreatedObj.Phone,
			},
			newChat.Members[1]: {
				Id:    userReceive.Id,
				Nama:  userReceive.Nama,
				Phone: userReceive.Phone,
			},
		},
		Messages:  newChat.Messages,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		IsDeleted: false,
		Type:      1,
	}

	err := cs.chatRepo.CreateNewChat(newChatInsert)

	return newChatInsert, err

}

func (cs chatService) SendChat(chat domain.ChatBubble, chatRoomId string) error {

	return cs.chatRepo.SendChat(chat, chatRoomId)
	// return nil

}

func (cs chatService) SetReadedChat(chatRoomId, chatBubbleId string) error {

	err := cs.chatRepo.SetReadedChat(chatRoomId, chatBubbleId)

	return err

}

func (cs chatService) GetChatRoomId(chatRoomId string) (*domain.Chat, error) {

	return cs.chatRepo.GetChatRoomId(chatRoomId)

}

func (cs chatService) GetChatRooms(userId string) ([]*domain.Chat, error) {

	return cs.chatRepo.GetChatRooms(userId)

}
