package main

import (
	"fmt"
	"library-api/internal/config"
	"library-api/internal/domain"
	"library-api/pkg/database"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg := config.Load()
	fmt.Println("CFG: ", cfg)

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create admin user
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := domain.User{
		Username: "admin",
		Email:    "admin@library.com",
		Password: string(adminPassword),
		Role:     "admin",
	}

	if err := db.FirstOrCreate(&admin, domain.User{Username: "admin"}).Error; err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	// Create regular user
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := domain.User{
		Username: "user",
		Email:    "user@library.com",
		Password: string(userPassword),
		Role:     "user",
	}

	if err := db.FirstOrCreate(&user, domain.User{Username: "user"}).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}

	// Create sample books
	books := []domain.Book{
		{
			Title:       "The Go Programming Language",
			Author:      "Alan A. A. Donovan, Brian W. Kernighan",
			ISBN:        "978-0134190440",
			Description: "A comprehensive guide to Go programming",
			Stock:       5,
			Available:   5,
		},
		{
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			ISBN:        "978-0132350884",
			Description: "A handbook of agile software craftsmanship",
			Stock:       3,
			Available:   3,
		},
		{
			Title:       "Design Patterns",
			Author:      "Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides",
			ISBN:        "978-0201633610",
			Description: "Elements of reusable object-oriented software",
			Stock:       2,
			Available:   2,
		},
	}

	for _, book := range books {
		if err := db.FirstOrCreate(&book, domain.Book{ISBN: book.ISBN}).Error; err != nil {
			log.Printf("Failed to create book %s: %v", book.Title, err)
		}
	}

	log.Println("Database seeded successfully!")
	log.Println("Admin credentials: admin/admin123")
	log.Println("User credentials: user/user123")
}
