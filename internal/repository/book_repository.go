package repository

import (
	"library-api/internal/domain"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *domain.Book) error
	GetByID(id uint) (*domain.Book, error)
	GetByISBN(isbn string) (*domain.Book, error)
	Update(book *domain.Book) error
	Delete(id uint) error
	GetAll() ([]*domain.Book, error)
	GetAvailable() ([]*domain.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *domain.Book) error {
	book.Available = book.Stock
	return r.db.Create(book).Error
}

func (r *bookRepository) GetByID(id uint) (*domain.Book, error) {
	var book domain.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetByISBN(isbn string) (*domain.Book, error) {
	var book domain.Book
	err := r.db.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Book{}, id).Error
}

func (r *bookRepository) GetAll() ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Find(&books).Error
	return books, err
}

func (r *bookRepository) GetAvailable() ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Where("available > 0").Find(&books).Error
	return books, err
}
