package middleware

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		// Map routes to resources and actions
		resource, action := mapRouteToResourceAction(c.Request.Method, c.FullPath())

		// Check permission
		allowed, err := enforcer.Enforce(role, resource, action)
		fmt.Println("Checking permission for role:", role, "resource:", resource, "action:", action, "allowed:", allowed)
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
	default:
		return "unknown", "unknown"
	}
}
