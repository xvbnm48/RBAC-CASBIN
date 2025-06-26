package main

import (
	"fmt"
	"library-api/internal/config"
	"library-api/internal/domain"
	"library-api/internal/repository"
	"library-api/pkg/database"
	"library-api/pkg/rbac"
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

	// Initialize repositories for RBAC setup
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)

	// Initialize RBAC service
	rbacService, err := rbac.NewRBACService(db, userRoleRepo, roleRepo, permissionRepo)
	if err != nil {
		log.Fatal("Failed to initialize RBAC service:", err)
	}

	// Setup default policies and roles
	if err := rbacService.SetupDefaultPolicies(); err != nil {
		log.Fatal("Failed to setup default policies:", err)
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

	// Assign roles to users (idempotent)
	adminRole, err := roleRepo.GetByName("admin")
	if err != nil {
		log.Fatal("Failed to get admin role:", err)
	}

	userRole, err := roleRepo.GetByName("user")
	if err != nil {
		log.Fatal("Failed to get user role:", err)
	}

	// Check if admin user already has admin role
	existingAdminRoles, err := userRoleRepo.GetUserRoles(admin.ID)
	if err == nil && len(existingAdminRoles) == 0 {
		// Assign admin role to admin user only if not already assigned
		if err := userRoleRepo.AssignRoles(admin.ID, []uint{adminRole.ID}); err != nil {
			log.Printf("Failed to assign admin role: %v", err)
		}
	}

	// Check if regular user already has user role
	existingUserRoles, err := userRoleRepo.GetUserRoles(user.ID)
	if err == nil && len(existingUserRoles) == 0 {
		// Assign user role to regular user only if not already assigned
		if err := userRoleRepo.AssignRoles(user.ID, []uint{userRole.ID}); err != nil {
			log.Printf("Failed to assign user role: %v", err)
		}
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
