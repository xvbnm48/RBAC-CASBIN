package middleware

import (
	"fmt"
	"library-api/pkg/rbac"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(rbacService *rbac.RBACService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			c.Abort()
			return
		}

		// Map routes to resources and actions
		resource, action := mapRouteToResourceAction(c.Request.Method, c.FullPath())
		fmt.Println("Resource:", resource, "Action:", action)

		// Check permission using new RBAC service
		allowed, err := rbacService.CheckPermission(userID.(uint), resource, action)
		fmt.Println("Permission check result:", allowed, "Error:", err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func mapRouteToResourceAction(method, path string) (resource, action string) {
	switch {
	case path == "/api/books" && method == "GET":
		return "books", "read"
	case path == "/api/books/:id" && method == "GET":
		return "books", "read"
	case path == "/api/books" && method == "POST":
		return "books", "write"
	case path == "/api/books/:id" && method == "PUT":
		return "books", "write"
	case path == "/api/books/:id" && method == "DELETE":
		return "books", "delete"
	case path == "/api/borrowings" && method == "GET":
		return "borrowings", "read"
	case path == "/api/borrowings" && method == "POST":
		return "borrowings", "write"
	case path == "/api/borrowings/:id/return" && method == "PUT":
		return "borrowings", "write"
	case path == "/api/users" && method == "GET":
		return "users", "read"
	case path == "/api/users/:id" && method == "GET":
		return "users", "read"
	case path == "/api/users/:id" && method == "PUT":
		return "users", "write"
	// Role management endpoints
	case path == "/api/roles" && method == "GET":
		return "roles", "read"
	case path == "/api/roles/:id" && method == "GET":
		return "roles", "read"
	case path == "/api/roles" && method == "POST":
		return "roles", "write"
	case path == "/api/roles/:id" && method == "PUT":
		return "roles", "write"
	case path == "/api/roles/:id" && method == "DELETE":
		return "roles", "delete"
	case path == "/api/roles/assign" && method == "POST":
		return "roles", "write"
	case path == "/api/users/:id/roles" && method == "GET":
		return "users", "read"
	// Permission management endpoints
	case path == "/api/permissions" && method == "GET":
		return "permissions", "read"
	case path == "/api/permissions/:id" && method == "GET":
		return "permissions", "read"
	case path == "/api/permissions" && method == "POST":
		return "permissions", "write"
	case path == "/api/permissions/:id" && method == "DELETE":
		return "permissions", "delete"
	default:
		return "unknown", "unknown"
	}
}
