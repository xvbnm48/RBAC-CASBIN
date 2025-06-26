package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;default:'user'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Borrowings []Borrowing `json:"borrowings,omitempty" gorm:"foreignKey:UserID"`
	UserRoles  []UserRole  `json:"user_roles,omitempty" gorm:"foreignKey:UserID"`
}

type Book struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Author      string         `json:"author" gorm:"not null"`
	ISBN        string         `json:"isbn" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Stock       int            `json:"stock" gorm:"not null;default:0"`
	Available   int            `json:"available" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Borrowings []Borrowing `json:"borrowings,omitempty" gorm:"foreignKey:BookID"`
}

type Borrowing struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	BookID     uint           `json:"book_id" gorm:"not null"`
	BorrowDate time.Time      `json:"borrow_date" gorm:"not null"`
	DueDate    time.Time      `json:"due_date" gorm:"not null"`
	ReturnDate *time.Time     `json:"return_date"`
	Status     string         `json:"status" gorm:"not null;default:'borrowed'"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
	Book Book `json:"book" gorm:"foreignKey:BookID"`
}

// Request/Response DTOs
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	ISBN        string `json:"isbn" binding:"required"`
	Description string `json:"description"`
	Stock       int    `json:"stock" binding:"required,min=0"`
}

type UpdateBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"`
	Stock       int    `json:"stock,omitempty"`
}

type BorrowBookRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

type BorrowingResponse struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	BookID     uint       `json:"book_id"`
	BookTitle  string     `json:"book_title"`
	Author     string     `json:"author"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
	Status     string     `json:"status"`
}

// Role represents a role in the system
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	UserRoles   []UserRole   `json:"user_roles,omitempty" gorm:"foreignKey:RoleID"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// Permission represents a permission in the system
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Resource    string         `json:"resource" gorm:"not null"`
	Action      string         `json:"action" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	RoleID    uint           `json:"role_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User User `json:"user" gorm:"foreignKey:UserID"`
	Role Role `json:"role" gorm:"foreignKey:RoleID"`
}

// Request/Response DTOs for Role Management
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions"` // Array of permission IDs
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions"` // Array of permission IDs
}

type AssignRoleRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Roles  []uint `json:"roles" binding:"required"` // Array of role IDs
}

type CreatePermissionRequest struct {
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
}

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type UserWithRolesResponse struct {
	ID       uint            `json:"id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Roles    []*RoleResponse `json:"roles"`
}
