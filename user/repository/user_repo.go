package repository

import (
	"adzi-clean-architecture/domain"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	userCollection *mongo.Collection
}

func NewUserRepo(collection *mongo.Collection) domain.UserRepository {

	return &userRepo{
		userCollection: collection,
	}
}

func (userRepo *userRepo) GetUserAll() []*domain.User {

	filter := bson.D{}

	var users []*domain.User

	data, err := userRepo.userCollection.Find(context.Background(), filter)

	if err != nil {
		panic(err)
	}

	if err := data.All(context.Background(), &users); err != nil {
		log.Fatal(err)
	}

	return users

}

func (userRepo *userRepo) GetUserById(id string) (*domain.User, error) {

	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objId}

	var user *domain.User

	projection := bson.M{
		"_id":   1,
		"nama":  1,
		"phone": 1,
	}

	err := userRepo.userCollection.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&user)

	return user, err

}

func (userRepo *userRepo) GetUserByEmail(email string) (*domain.User, error) {

	// objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"email": email}

	var user *domain.User

	err := userRepo.userCollection.FindOne(context.Background(), filter).Decode(&user)

	// fmt.Println("email foorm repo>>>>>>>>>.", email)
	// fmt.Println("user foorm repo>>>>>>>>>.", user)

	return user, err

}
