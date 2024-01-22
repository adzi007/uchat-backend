package service

import (
	"adzi-clean-architecture/domain"
	"adzi-clean-architecture/domain/mocks"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserAll(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	id1, _ := primitive.ObjectIDFromHex("62cda02bf5a4d0c8cbf3795a")

	mockUsers := domain.User{
		Id:           id1,
		Nama:         "Udin",
		Username:     "udin2023",
		Password:     "pas123",
		Email:        "udinxxxx2023@gmail.com",
		Phone:        "09899887",
		Token:        "",
		UserType:     "",
		RefreshToken: "",
	}

	mockListUser := make([]*domain.User, 0)
	mockListUser = append(mockListUser, &mockUsers)

	// mockUserRepo.On("GetUserAll")

	mockUserRepo.On("GetUserAll").Return(mockListUser)

	// Arrange
	realService := NewUserService(mockUserRepo)
	// mockService := new(mocks.UserService)

	resultMocked := realService.GetUserAll()

	fmt.Println(resultMocked[0].Nama)

	assert.Equal(t, mockUsers.Id, resultMocked[0].Id)
	assert.Equal(t, mockUsers.Nama, resultMocked[0].Nama)
	assert.Equal(t, mockUsers.Username, resultMocked[0].Username)
	assert.Equal(t, mockUsers.Password, resultMocked[0].Password)

}
