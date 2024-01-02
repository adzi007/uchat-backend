package repository

import (
	"adzi-clean-architecture/domain"
	"context"
	"time"

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

func (cr *chatRepo) SetReadedChat(chatRoomId, chatBubbleId string) error {

	objChatRoomId, _ := primitive.ObjectIDFromHex(chatRoomId)
	objChatBubbleId, _ := primitive.ObjectIDFromHex(chatBubbleId)

	filter := bson.M{
		"_id":          objChatRoomId,
		"messages._id": objChatBubbleId,
	}

	update := bson.M{
		"$set": bson.M{
			"messages.$.readedat": time.Now().UTC(),
		},
	}
	_, err := cr.chatCollection.UpdateOne(context.Background(), filter, update)

	return err
}

func (cr *chatRepo) GetChatRoomId(chatRoomId string) (*domain.Chat, error) {

	objChatRoomId, _ := primitive.ObjectIDFromHex(chatRoomId)

	filter := bson.M{
		"_id": objChatRoomId,
	}

	var chatRoom *domain.Chat

	err := cr.chatCollection.FindOne(context.Background(), filter).Decode(&chatRoom)

	// fmt.Println("chat repo res", chatRoom)
	// fmt.Println("chat repo err", err)

	return chatRoom, err

}
