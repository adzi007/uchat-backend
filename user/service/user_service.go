package service

import (
	"adzi-clean-architecture/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(ur domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: ur,
	}
}

func (us userService) GetUserAll() []*domain.User {

	return us.userRepo.GetUserAll()

}

func (us *userService) GetUserById(id string) (*domain.User, error) {

	data, err := us.userRepo.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return data, nil

}
