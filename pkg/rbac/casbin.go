package rbac

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	// Initialize GORM adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// Initialize Casbin enforcer
	enforcer, err := casbin.NewEnforcer("configs/casbin_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// Load policies from file (optional, for initial setup)
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

func SetupDefaultPolicies(enforcer *casbin.Enforcer) error {
	// Add policies for admin role
	policies := [][]string{
		{"admin", "books", "read"},
		{"admin", "books", "write"},
		{"admin", "books", "delete"},
		{"admin", "borrowings", "read"},
		{"admin", "borrowings", "write"},
		{"admin", "users", "read"},
		{"admin", "users", "write"},
		// Add policies for user role
		{"user", "books", "read"},
		{"user", "borrowings", "read"},
		{"user", "borrowings", "write"},
	}

	for _, policy := range policies {
		_, err := enforcer.AddPolicy(policy)
		if err != nil {
			return err
		}
	}

	return enforcer.SavePolicy()
}
