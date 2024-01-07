package repository

import (
	"adzi-clean-architecture/domain"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	return chatRoom, err

}

func (cr *chatRepo) GetChatRooms(userId string) ([]*domain.Chat, error) {

	objUserId, _ := primitive.ObjectIDFromHex(userId)

	filter := bson.M{
		"members": bson.M{
			"$elemMatch": bson.M{"$eq": objUserId},
		},
	}

	var chatRoom []*domain.Chat

	projection := bson.M{
		"_id":         1,
		"timeCreated": 1,
		"messages": bson.M{
			"$slice": bson.A{"$messages", -1},
		},
	}

	data, err := cr.chatCollection.Find(context.Background(), filter, options.Find().SetProjection(projection))

	if err != nil {

		fmt.Println("error disini", err.Error())
		panic(err)
	}

	if err := data.All(context.Background(), &chatRoom); err != nil {
		fmt.Println("error disini", err.Error())
		log.Fatal(err)
	}

	return chatRoom, err

}
