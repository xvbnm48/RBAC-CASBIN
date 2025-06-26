package service

import (
	"errors"
	"library-api/internal/domain"
	"library-api/internal/repository"
	"library-api/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *domain.RegisterRequest) (*domain.User, error)
	Login(req *domain.LoginRequest) (*domain.AuthResponse, error)
	// Test() utils.SuccessResponse // Removed or update to correct return type if needed
}

type authService struct {
	userRepo repository.UserRepository
	jwtAuth  *auth.JWTAuth
}

func NewAuthService(userRepo repository.UserRepository, jwtAuth *auth.JWTAuth) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (s *authService) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// Check if username already exists
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Validate role
	if role != "admin" && role != "user" {
		return nil, errors.New("invalid role")
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtAuth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
