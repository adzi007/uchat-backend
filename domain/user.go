package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Nama         string             `json:"nama"`
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	Token        string             `json:"token"`
	UserType     string             `json:"user_type"`
	RefreshToken string             `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserID       string             `json:"user_id"`
}

type UserService interface {
	GetUserAll() []*User
	GetUserById(string) (*User, error)
}

type UserRepository interface {
	GetUserAll() []*User
	GetUserById(string) (*User, error)
	GetUserByEmail(string) (*User, error)
}
