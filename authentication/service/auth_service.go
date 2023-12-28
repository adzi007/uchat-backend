package service

import (
	"adzi-clean-architecture/domain"
	"adzi-clean-architecture/pkg/logger"
	"adzi-clean-architecture/pkg/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo domain.UserRepository
}

func NewauthService(ur domain.UserRepository) domain.AuthService {
	return &authService{
		userRepo: ur,
	}
}

func (us authService) Login(loginRequest *domain.LoginRequest) (string, int, error) {

	// CHECK AVAILABLE USER

	getUser, err := us.userRepo.GetUserByEmail(loginRequest.Email)

	if err != nil {
		logger.Error().Err(err)
	}

	if getUser == nil {
		return "", 404, errors.New("user not found")
	}

	// CHECK PASSWORD

	isValid := utils.CheckPasswordHash(loginRequest.Password, getUser.Password)

	if !isValid {

		return "", 401, errors.New("wrong credential")
	}

	// GENERATE JWT

	claims := jwt.MapClaims{}

	claims["name"] = getUser.Nama
	claims["email"] = getUser.Email
	claims["exp"] = time.Now().Add(time.Minute * 3).Unix()

	if getUser.Email == "naruto@zmail.com" {
		claims["role"] = "admin"
	} else {
		claims["role"] = "user"
	}

	token, errGenerateToken := utils.GenerateToken(&claims)

	if errGenerateToken != nil {
		return "", 500, errGenerateToken
	}

	return token, 200, nil

}
