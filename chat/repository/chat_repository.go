package repository

import (
	"adzi-clean-architecture/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type chatRepo struct {
	chatCollection *mongo.Collection
}

func NewChatRepo(collection *mongo.Collection) domain.ChatRepository {

	return &chatRepo{
		chatCollection: collection,
	}
}

func (cr *chatRepo) CreateNewChat(newChat domain.Chat) error {

	_, err := cr.chatCollection.InsertOne(context.Background(), newChat)

	return err

}

func (cr *chatRepo) SendChat(chat domain.ChatBubble, chatRoomId string) error {

	objId, _ := primitive.ObjectIDFromHex(chatRoomId)
	filter := bson.M{"_id": objId}

	update := bson.M{
		"$push": bson.M{
			"messages": chat,
		},
	}

	_, err := cr.chatCollection.UpdateOne(context.Background(), filter, update)

	return err
}
