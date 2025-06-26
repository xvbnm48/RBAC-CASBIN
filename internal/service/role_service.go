package service

import (
	"errors"
	"fmt"
	"library-api/internal/domain"
	"library-api/internal/repository"

	"gorm.io/gorm"
)

type RoleService interface {
	CreateRole(req *domain.CreateRoleRequest) (*domain.RoleResponse, error)
	GetRoleByID(id uint) (*domain.RoleResponse, error)
	UpdateRole(id uint, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error)
	DeleteRole(id uint) error
	GetAllRoles() ([]*domain.RoleResponse, error)
	AssignRolesToUser(req *domain.AssignRoleRequest) error
	GetUserRoles(userID uint) ([]*domain.RoleResponse, error)
	GetUserWithRoles(userID uint) (*domain.UserWithRolesResponse, error)
}

type roleService struct {
	roleRepo     repository.RoleRepository
	userRepo     repository.UserRepository
	userRoleRepo repository.UserRoleRepository
	permRepo     repository.PermissionRepository
}

func NewRoleService(
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	permRepo repository.PermissionRepository,
) RoleService {
	return &roleService{
		roleRepo:     roleRepo,
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		permRepo:     permRepo,
	}
}

func (s *roleService) CreateRole(req *domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	// Check if role name already exists
	if _, err := s.roleRepo.GetByName(req.Name); err == nil {
		return nil, errors.New("role with this name already exists")
	}

	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	// Assign permissions if provided
	if len(req.Permissions) > 0 {
		if err := s.roleRepo.AssignPermissions(role.ID, req.Permissions); err != nil {
			return nil, err
		}
	}

	// Get role with permissions
	roleWithPerms, err := s.roleRepo.GetWithPermissions(role.ID)
	if err != nil {
		return nil, err
	}

	return s.roleToResponse(roleWithPerms), nil
}

func (s *roleService) GetRoleByID(id uint) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.GetWithPermissions(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return s.roleToResponse(role), nil
}

func (s *roleService) UpdateRole(id uint, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if new name already exists
		if existingRole, err := s.roleRepo.GetByName(req.Name); err == nil && existingRole.ID != role.ID {
			return nil, errors.New("role with this name already exists")
		}
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	// Update permissions if provided
	if req.Permissions != nil {
		if err := s.roleRepo.AssignPermissions(role.ID, req.Permissions); err != nil {
			return nil, err
		}
	}

	// Get updated role with permissions
	roleWithPerms, err := s.roleRepo.GetWithPermissions(role.ID)
	if err != nil {
		return nil, err
	}

	return s.roleToResponse(roleWithPerms), nil
}

func (s *roleService) DeleteRole(id uint) error {
	if _, err := s.roleRepo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	return s.roleRepo.Delete(id)
}

func (s *roleService) GetAllRoles() ([]*domain.RoleResponse, error) {
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*domain.RoleResponse
	for _, role := range roles {
		responses = append(responses, s.roleToResponse(role))
	}

	return responses, nil
}

func (s *roleService) AssignRolesToUser(req *domain.AssignRoleRequest) error {
	// Check if user exists
	fmt.Println("AssignRolesToUser called, param req:", req)
	if _, err := s.userRepo.GetByID(req.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Validate all role IDs exist
	for _, roleID := range req.Roles {
		if _, err := s.roleRepo.GetByID(roleID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("one or more roles not found")
			}
			return err
		}
	}

	return s.userRoleRepo.AssignRoles(req.UserID, req.Roles)
}

func (s *roleService) GetUserRoles(userID uint) ([]*domain.RoleResponse, error) {
	roles, err := s.userRoleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	var responses []*domain.RoleResponse
	for _, role := range roles {
		responses = append(responses, s.roleToResponse(role))
	}

	return responses, nil
}

func (s *roleService) GetUserWithRoles(userID uint) (*domain.UserWithRolesResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithRolesResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}, nil
}

func (s *roleService) roleToResponse(role *domain.Role) *domain.RoleResponse {
	var permissions []domain.PermissionResponse
	for _, perm := range role.Permissions {
		permissions = append(permissions, domain.PermissionResponse{
			ID:          perm.ID,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
		})
	}

	return &domain.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
	}
}
