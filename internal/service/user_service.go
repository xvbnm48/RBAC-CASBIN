package service

import (
	"library-api/internal/domain"
	"library-api/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]*domain.User, error)
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

type userService struct {
	userRepo repository.UserRepository
}

func (s *userService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAll()
}
