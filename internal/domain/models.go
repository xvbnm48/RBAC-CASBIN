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
