package service

import "adzi-clean-architecture/domain"

type chatService struct {
	chatRepo domain.ChatRepository
}

func NewChatService(cr domain.ChatRepository) domain.ChatService {
	return &chatService{
		chatRepo: cr,
	}
}

func (cs chatService) CreateNewChat(newChat domain.Chat) error {

	err := cs.chatRepo.CreateNewChat(newChat)

	return err

}

func (cs chatService) SendChat(chat domain.ChatBubble, chatRoomId string) error {

	// return cs.chatRepo.SendChat(chat, chatRoomId)
	return nil

}

func (cs chatService) SetReadedChat(chatRoomId, chatBubbleId string) error {

	err := cs.chatRepo.SetReadedChat(chatRoomId, chatBubbleId)

	return err

}

func (cs chatService) GetChatRoomId(chatRoomId string) (*domain.Chat, error) {

	return cs.chatRepo.GetChatRoomId(chatRoomId)

}
