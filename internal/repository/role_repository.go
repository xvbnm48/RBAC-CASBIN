package repository

import (
	"library-api/internal/domain"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *domain.Role) error
	GetByID(id uint) (*domain.Role, error)
	GetByName(name string) (*domain.Role, error)
	Update(role *domain.Role) error
	Delete(id uint) error
	GetAll() ([]*domain.Role, error)
	GetWithPermissions(id uint) (*domain.Role, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	ClearPermissions(roleID uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id uint) (*domain.Role, error) {
	var role domain.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(role *domain.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Role{}, id).Error
}

func (r *roleRepository) GetAll() ([]*domain.Role, error) {
	var roles []*domain.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) GetWithPermissions(id uint) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	role, err := r.GetByID(roleID)
	if err != nil {
		return err
	}

	// Clear existing permissions
	err = r.db.Model(role).Association("Permissions").Clear()
	if err != nil {
		return err
	}

	// Add new permissions
	if len(permissionIDs) > 0 {
		var permissions []domain.Permission
		err = r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error
		if err != nil {
			return err
		}

		err = r.db.Model(role).Association("Permissions").Append(permissions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *roleRepository) ClearPermissions(roleID uint) error {
	role, err := r.GetByID(roleID)
	if err != nil {
		return err
	}

	// Clear all permissions for this role
	return r.db.Model(role).Association("Permissions").Clear()
}
