package repository

import (
	"library-api/internal/domain"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *domain.Permission) error
	GetByID(id uint) (*domain.Permission, error)
	GetByResourceAndAction(resource, action string) (*domain.Permission, error)
	Update(permission *domain.Permission) error
	Delete(id uint) error
	GetAll() ([]*domain.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *domain.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetByID(id uint) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByResourceAndAction(resource, action string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(permission *domain.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Permission{}, id).Error
}

func (r *permissionRepository) GetAll() ([]*domain.Permission, error) {
	var permissions []*domain.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}
