package handler

import (
	"library-api/internal/domain"
	"library-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BorrowingHandler struct {
	borrowingService service.BorrowingService
}

func NewBorrowingHandler(borrowingService service.BorrowingService) *BorrowingHandler {
	return &BorrowingHandler{
		borrowingService: borrowingService,
	}
}

func (h *BorrowingHandler) BorrowBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var req domain.BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrowing, err := h.borrowingService.BorrowBook(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, borrowing)
}

func (h *BorrowingHandler) ReturnBook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	idStr := c.Param("id")
	borrowingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrowing ID"})
		return
	}

	borrowing, err := h.borrowingService.ReturnBook(uint(borrowingID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrowing)
}

func (h *BorrowingHandler) GetBorrowings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
		return
	}

	var borrowings []*domain.BorrowingResponse
	var err error

	// Admin can see all borrowings, user can only see their own
	if role.(string) == "admin" {
		borrowings, err = h.borrowingService.GetAllBorrowings()
	} else {
		borrowings, err = h.borrowingService.GetUserBorrowings(userID.(uint))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrowings)
}
