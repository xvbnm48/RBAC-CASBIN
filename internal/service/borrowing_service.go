package service

import (
	"errors"
	"library-api/internal/domain"
	"library-api/internal/repository"
	"time"

	"gorm.io/gorm"
)

type BorrowingService interface {
	BorrowBook(userID uint, req *domain.BorrowBookRequest) (*domain.BorrowingResponse, error)
	ReturnBook(borrowingID uint, userID uint) (*domain.BorrowingResponse, error)
	GetUserBorrowings(userID uint) ([]*domain.BorrowingResponse, error)
	GetAllBorrowings() ([]*domain.BorrowingResponse, error)
}

type borrowingService struct {
	borrowingRepo repository.BorrowingRepository
	bookRepo      repository.BookRepository
}

func NewBorrowingService(borrowingRepo repository.BorrowingRepository, bookRepo repository.BookRepository) BorrowingService {
	return &borrowingService{
		borrowingRepo: borrowingRepo,
		bookRepo:      bookRepo,
	}
}

func (s *borrowingService) BorrowBook(userID uint, req *domain.BorrowBookRequest) (*domain.BorrowingResponse, error) {
	// Check if book exists and is available
	book, err := s.bookRepo.GetByID(req.BookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if book.Available <= 0 {
		return nil, errors.New("book is not available")
	}

	// Check if user already borrowed this book and hasn't returned it
	if _, err := s.borrowingRepo.GetActiveBorrowingByUserAndBook(userID, req.BookID); err == nil {
		return nil, errors.New("you have already borrowed this book")
	}

	// Create borrowing record
	borrowing := &domain.Borrowing{
		UserID: userID,
		BookID: req.BookID,
	}

	if err := s.borrowingRepo.Create(borrowing); err != nil {
		return nil, err
	}

	// Update book availability
	book.Available--
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	// Get the created borrowing with relations
	createdBorrowing, err := s.borrowingRepo.GetByID(borrowing.ID)
	if err != nil {
		return nil, err
	}

	return s.borrowingToResponse(createdBorrowing), nil
}

func (s *borrowingService) ReturnBook(borrowingID uint, userID uint) (*domain.BorrowingResponse, error) {
	borrowing, err := s.borrowingRepo.GetByID(borrowingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("borrowing record not found")
		}
		return nil, err
	}

	// Check if the borrowing belongs to the user
	if borrowing.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Check if already returned
	if borrowing.Status == "returned" {
		return nil, errors.New("book already returned")
	}

	// Update borrowing record
	now := time.Now()
	borrowing.ReturnDate = &now
	borrowing.Status = "returned"

	if err := s.borrowingRepo.Update(borrowing); err != nil {
		return nil, err
	}

	// Update book availability
	book, err := s.bookRepo.GetByID(borrowing.BookID)
	if err != nil {
		return nil, err
	}

	book.Available++
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return s.borrowingToResponse(borrowing), nil
}

func (s *borrowingService) GetUserBorrowings(userID uint) ([]*domain.BorrowingResponse, error) {
	borrowings, err := s.borrowingRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []*domain.BorrowingResponse
	for _, borrowing := range borrowings {
		responses = append(responses, s.borrowingToResponse(borrowing))
	}

	return responses, nil
}

func (s *borrowingService) GetAllBorrowings() ([]*domain.BorrowingResponse, error) {
	borrowings, err := s.borrowingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*domain.BorrowingResponse
	for _, borrowing := range borrowings {
		responses = append(responses, s.borrowingToResponse(borrowing))
	}

	return responses, nil
}

func (s *borrowingService) borrowingToResponse(borrowing *domain.Borrowing) *domain.BorrowingResponse {
	return &domain.BorrowingResponse{
		ID:         borrowing.ID,
		UserID:     borrowing.UserID,
		BookID:     borrowing.BookID,
		BookTitle:  borrowing.Book.Title,
		Author:     borrowing.Book.Author,
		BorrowDate: borrowing.BorrowDate,
		DueDate:    borrowing.DueDate,
		ReturnDate: borrowing.ReturnDate,
		Status:     borrowing.Status,
	}
}
