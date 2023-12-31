package domain

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewChatRequest struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	TimeCreated time.Time          `json:"timeCreated"`
	Members     []string           `json:"members"`
	// Members   map[string]User `json:"members"`
	Messages []ChatBubble `json:"chatBubble"`
}

type Chat struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	TimeCreated time.Time          `json:"timeCreated"`
	// Members     []primitive.ObjectID `json:"members"`
	Members   map[string]*UserResponseData `json:"members"`
	Messages  []ChatBubble                 `json:"chatBubble"`
	Type      int32                        `json:"type"`
	CreatedAt time.Time                    `json:"created_at"`
	UpdatedAt time.Time                    `json:"updated_at"`
	IsDeleted bool                         `json:"isDeleted"`
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
	CreateNewChat(NewChatRequest) (Chat, error)
	SendChat(ChatBubble, string) error
	SetReadedChat(chatRoomId, chatBubbleId string) error
	GetChatRoomId(chatRoomId string) (*Chat, error)
	GetChatRooms(userId string) ([]*Chat, error)
}

type ChatRepository interface {
	CreateNewChat(Chat) error
	SendChat(ChatBubble, string) error
	SetReadedChat(chatRoomId, chatBubbleId string) error
	GetChatRoomId(chatRoomId string) (*Chat, error)
	GetChatRooms(userId string) ([]*Chat, error)
}

type ChatWebsocket interface {
	Run()
	Join(*websocket.Conn, string)
	Leave(*websocket.Conn)
	Broadcast(ChatBubble, string)
	HandleWsChatRoom() func(*websocket.Conn)
}
