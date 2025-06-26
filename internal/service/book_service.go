package service

import (
	"errors"
	"library-api/internal/domain"
	"library-api/internal/repository"

	"gorm.io/gorm"
)

type BookService interface {
	Create(req *domain.CreateBookRequest) (*domain.Book, error)
	GetByID(id uint) (*domain.Book, error)
	Update(id uint, req *domain.UpdateBookRequest) (*domain.Book, error)
	Delete(id uint) error
	GetAll() ([]*domain.Book, error)
	GetAvailable() ([]*domain.Book, error)
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (s *bookService) Create(req *domain.CreateBookRequest) (*domain.Book, error) {
	// Check if ISBN already exists
	if _, err := s.bookRepo.GetByISBN(req.ISBN); err == nil {
		return nil, errors.New("book with this ISBN already exists")
	}

	book := &domain.Book{
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		Description: req.Description,
		Stock:       req.Stock,
		Available:   req.Stock,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) GetByID(id uint) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return book, nil
}

func (s *bookService) Update(id uint, req *domain.UpdateBookRequest) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.ISBN != "" {
		// Check if new ISBN already exists
		if existingBook, err := s.bookRepo.GetByISBN(req.ISBN); err == nil && existingBook.ID != book.ID {
			return nil, errors.New("book with this ISBN already exists")
		}
		book.ISBN = req.ISBN
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Stock > 0 {
		// Calculate the difference and update available count
		stockDiff := req.Stock - book.Stock
		book.Stock = req.Stock
		book.Available += stockDiff
		if book.Available < 0 {
			book.Available = 0
		}
	}

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) Delete(id uint) error {
	if _, err := s.bookRepo.GetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("book not found")
		}
		return err
	}

	return s.bookRepo.Delete(id)
}

func (s *bookService) GetAll() ([]*domain.Book, error) {
	return s.bookRepo.GetAll()
}

func (s *bookService) GetAvailable() ([]*domain.Book, error) {
	return s.bookRepo.GetAvailable()
}
