package main

import (
	"library-api/internal/config"
	"library-api/internal/handler"
	"library-api/internal/middleware"
	"library-api/internal/repository"
	"library-api/internal/service"
	"library-api/pkg/auth"
	"library-api/pkg/database"
	"library-api/pkg/rbac"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	borrowingRepo := repository.NewBorrowingRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)

	// Initialize RBAC service
	rbacService, err := rbac.NewRBACService(db, userRoleRepo, roleRepo, permissionRepo)
	if err != nil {
		log.Fatal("Failed to initialize RBAC service:", err)
	}

	// Setup default policies
	if err := rbacService.SetupDefaultPolicies(); err != nil {
		log.Fatal("Failed to setup default policies:", err)
	}

	// Initialize JWT auth
	jwtAuth := auth.NewJWTAuth(cfg.JWT.SecretKey, cfg.JWT.ExpiryHours)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtAuth)
	bookService := service.NewBookService(bookRepo)
	borrowingService := service.NewBorrowingService(borrowingRepo, bookRepo)
	roleService := service.NewRoleService(roleRepo, userRepo, userRoleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	bookHandler := handler.NewBookHandler(bookService)
	borrowingHandler := handler.NewBorrowingHandler(borrowingService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// Setup router
	router := gin.Default()

	// health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtAuth))
	api.Use(middleware.RBACMiddleware(rbacService))
	{
		// Book routes
		books := api.Group("/books")
		{
			books.GET("", bookHandler.GetAll)
			books.GET("/:id", bookHandler.GetByID)
			books.POST("", bookHandler.Create)
			books.PUT("/:id", bookHandler.Update)
			books.DELETE("/:id", bookHandler.Delete)
		}

		// Borrowing routes
		borrowings := api.Group("/borrowings")
		{
			borrowings.GET("", borrowingHandler.GetBorrowings)
			borrowings.POST("", borrowingHandler.BorrowBook)
			borrowings.PUT("/:id/return", borrowingHandler.ReturnBook)
		}

		// Role management routes (Admin only)
		roles := api.Group("/roles")
		{
			roles.GET("", roleHandler.GetAllRoles)
			roles.GET("/:id", roleHandler.GetRoleByID)
			roles.POST("", roleHandler.CreateRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.POST("/assign", roleHandler.AssignRolesToUser)
		}

		// Permission management routes (Admin only)
		permissions := api.Group("/permissions")
		{
			permissions.GET("", permissionHandler.GetAllPermissions)
			permissions.GET("/:id", permissionHandler.GetPermissionByID)
			permissions.POST("", permissionHandler.CreatePermission)
			permissions.DELETE("/:id", permissionHandler.DeletePermission)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id/roles", roleHandler.GetUserRoles)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
