package repository

import (
	"fmt"
	"library-api/internal/domain"

	"gorm.io/gorm"
)

type UserRoleRepository interface {
	AssignRoles(userID uint, roleIDs []uint) error
	RemoveRoles(userID uint, roleIDs []uint) error
	GetUserRoles(userID uint) ([]*domain.Role, error)
	GetRoleUsers(roleID uint) ([]*domain.User, error)
	HasRole(userID uint, roleName string) (bool, error)
	ClearUserRoles(userID uint) error
}

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// First clear existing roles
	fmt.Println("Clearing existing roles for user ID:", userID)
	err := r.ClearUserRoles(userID)
	if err != nil {
		return err
	}
	fmt.Println("Existing roles cleared for user ID:", userID)

	// Add new roles
	for _, roleID := range roleIDs {
		userRole := &domain.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := r.db.Create(userRole).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRoleRepository) RemoveRoles(userID uint, roleIDs []uint) error {
	return r.db.Where("user_id = ? AND role_id IN ?", userID, roleIDs).Delete(&domain.UserRole{}).Error
}

func (r *userRoleRepository) GetUserRoles(userID uint) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.db.Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND user_roles.deleted_at IS NULL", userID).
		Preload("Permissions").
		Find(&roles).Error
	return roles, err
}

func (r *userRoleRepository) GetRoleUsers(roleID uint) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND user_roles.deleted_at IS NULL", roleID).
		Find(&users).Error
	return users, err
}

func (r *userRoleRepository) HasRole(userID uint, roleName string) (bool, error) {
	var count int64
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ? AND user_roles.deleted_at IS NULL", userID, roleName).
		Count(&count).Error
	return count > 0, err
}

func (r *userRoleRepository) ClearUserRoles(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserRole{}).Error
}
