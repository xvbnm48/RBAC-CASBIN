package service

import (
	"errors"
	"library-api/internal/domain"
	"library-api/internal/repository"

	"gorm.io/gorm"
)

type PermissionService interface {
	CreatePermission(req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error)
	GetPermissionByID(id uint) (*domain.PermissionResponse, error)
	GetAllPermissions() ([]*domain.PermissionResponse, error)
	DeletePermission(id uint) error
}

type permissionService struct {
	permRepo repository.PermissionRepository
}

func NewPermissionService(permRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		permRepo: permRepo,
	}
}

func (s *permissionService) CreatePermission(req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error) {
	// Check if permission already exists
	if _, err := s.permRepo.GetByResourceAndAction(req.Resource, req.Action); err == nil {
		return nil, errors.New("permission with this resource and action already exists")
	}

	permission := &domain.Permission{
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	if err := s.permRepo.Create(permission); err != nil {
		return nil, err
	}

	return s.permissionToResponse(permission), nil
}

func (s *permissionService) GetPermissionByID(id uint) (*domain.PermissionResponse, error) {
	permission, err := s.permRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return s.permissionToResponse(permission), nil
}

func (s *permissionService) GetAllPermissions() ([]*domain.PermissionResponse, error) {
	permissions, err := s.permRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*domain.PermissionResponse
	for _, permission := range permissions {
		responses = append(responses, s.permissionToResponse(permission))
	}

	return responses, nil
}

func (s *permissionService) DeletePermission(id uint) error {
	if _, err := s.permRepo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}
		return err
	}

	return s.permRepo.Delete(id)
}

func (s *permissionService) permissionToResponse(permission *domain.Permission) *domain.PermissionResponse {
	return &domain.PermissionResponse{
		ID:          permission.ID,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
	}
}
