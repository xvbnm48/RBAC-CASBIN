package repository

import (
	"library-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type BorrowingRepository interface {
	Create(borrowing *domain.Borrowing) error
	GetByID(id uint) (*domain.Borrowing, error)
	GetByUserID(userID uint) ([]*domain.Borrowing, error)
	GetAll() ([]*domain.Borrowing, error)
	Update(borrowing *domain.Borrowing) error
	GetActiveBorrowingByUserAndBook(userID, bookID uint) (*domain.Borrowing, error)
}

type borrowingRepository struct {
	db *gorm.DB
}

func NewBorrowingRepository(db *gorm.DB) BorrowingRepository {
	return &borrowingRepository{db: db}
}

func (r *borrowingRepository) Create(borrowing *domain.Borrowing) error {
	borrowing.BorrowDate = time.Now()
	borrowing.DueDate = time.Now().AddDate(0, 0, 14) // 14 hari peminjaman
	borrowing.Status = "borrowed"
	return r.db.Create(borrowing).Error
}

func (r *borrowingRepository) GetByID(id uint) (*domain.Borrowing, error) {
	var borrowing domain.Borrowing
	err := r.db.Preload("User").Preload("Book").First(&borrowing, id).Error
	if err != nil {
		return nil, err
	}
	return &borrowing, nil
}

func (r *borrowingRepository) GetByUserID(userID uint) ([]*domain.Borrowing, error) {
	var borrowings []*domain.Borrowing
	err := r.db.Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&borrowings).Error
	return borrowings, err
}

func (r *borrowingRepository) GetAll() ([]*domain.Borrowing, error) {
	var borrowings []*domain.Borrowing
	err := r.db.Preload("User").Preload("Book").Find(&borrowings).Error
	return borrowings, err
}

func (r *borrowingRepository) Update(borrowing *domain.Borrowing) error {
	return r.db.Save(borrowing).Error
}

func (r *borrowingRepository) GetActiveBorrowingByUserAndBook(userID, bookID uint) (*domain.Borrowing, error) {
	var borrowing domain.Borrowing
	err := r.db.Where("user_id = ? AND book_id = ? AND status = ?", userID, bookID, "borrowed").First(&borrowing).Error
	if err != nil {
		return nil, err
	}
	return &borrowing, nil
}
