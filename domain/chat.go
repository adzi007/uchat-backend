package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id"`
	TimeCreated time.Time            `json:"timeCreated"`
	Members     []primitive.ObjectID `json:"members"`
	Messages    []ChatBubble         `json:"chatBubble"`
	Type        string               `json:"type"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	IsDeleted   bool                 `json:"isDeleted"`
}

type DocumentAttachment struct {
	FileName string `json:"fileName"`
	Size     int32  `json:"size"`
}

type ChatBubbleMedia struct {
	Images   []string           `json:"images"`
	Video    string             `json:"video"`
	Document DocumentAttachment `json:"document"`
}

type ChatBubble struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Timestamp  time.Time          `json:"timestamp"`
	ReplyId    primitive.ObjectID `json:"replyId,omitempty" bson:"replyId,omitempty"`
	UserID     primitive.ObjectID `json:"userId"`
	Message    string             `json:"message"`
	Attachment ChatBubbleMedia    `json:"attachment"`
	IsDeleted  bool               `json:"isDeleted"`
	ReadedAt   *time.Time         `json:"readedAt,omitempty"`
}

type ChatCreateRequest struct {
	UserCreate   string `json:"userCreate" validate:"required"`
	UserReceiver string `json:"userReceiver" validate:"required"`
	Message      string `json:"message" validate:"required"`
	Attachment   string `json:"attachment"`
}

type CreateChatBubbleRequest struct {
	ChatRoomId string          `json:"chatRoomId" validate:"required"`
	ReplyId    string          `json:"replyId"`
	UserID     string          `json:"userId" validate:"required"`
	Message    string          `json:"message"`
	Attachment ChatBubbleMedia `json:"attachment"`
}

type SetReadedRequest struct {
	UserID       string `json:"userId" validate:"required"`
	ChatRoomId   string `json:"chatRoomId" validate:"required"`
	ChatBubbleId string `json:"chatBubbleId" validate:"required"`
}

type ChatService interface {
	CreateNewChat(Chat) error
	SendChat(ChatBubble, string) error
	SetReadedChat(chatRoomId, chatBubbleId string) error
}

type ChatRepository interface {
	CreateNewChat(Chat) error
	SendChat(ChatBubble, string) error
	SetReadedChat(chatRoomId, chatBubbleId string) error
}
