package rbac

import (
	"fmt"
	"library-api/internal/domain"
	"library-api/internal/repository"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type RBACService struct {
	enforcer     *casbin.Enforcer
	userRoleRepo repository.UserRoleRepository
	roleRepo     repository.RoleRepository
	permRepo     repository.PermissionRepository
}

func NewRBACService(
	db *gorm.DB,
	userRoleRepo repository.UserRoleRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
) (*RBACService, error) {
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

	// Load policies from database
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	service := &RBACService{
		enforcer:     enforcer,
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		permRepo:     permRepo,
	}

	return service, nil
}

func (s *RBACService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}

func (s *RBACService) CheckPermission(userID uint, resource, action string) (bool, error) {
	// Get user roles
	roles, err := s.userRoleRepo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	// Check each role for permission
	for _, role := range roles {
		allowed, err := s.enforcer.Enforce(role.Name, resource, action)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}

	return false, nil
}

func (s *RBACService) SyncRolePermissions() error {
	// Get all roles with permissions
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return err
	}

	// Add policies for each role, checking for existence first
	for _, role := range roles {
		for _, permission := range role.Permissions {
			// Check if policy already exists before adding
			if !s.enforcer.HasPolicy(role.Name, permission.Resource, permission.Action) {
				_, err := s.enforcer.AddPolicy(role.Name, permission.Resource, permission.Action)
				if err != nil {
					return fmt.Errorf("failed to add policy for role %s: %v", role.Name, err)
				}
			}
		}
	}

	// Save policies to database
	return s.enforcer.SavePolicy()
}

func (s *RBACService) AddRolePolicy(roleName, resource, action string) error {
	// Check if policy already exists
	if s.enforcer.HasPolicy(roleName, resource, action) {
		return nil // Policy already exists, no need to add
	}

	_, err := s.enforcer.AddPolicy(roleName, resource, action)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

func (s *RBACService) RemoveRolePolicy(roleName, resource, action string) error {
	_, err := s.enforcer.RemovePolicy(roleName, resource, action)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

func (s *RBACService) SetupDefaultPolicies() error {
	// Get or create default permissions
	defaultPermissions := []struct {
		Resource    string
		Action      string
		Description string
	}{
		{"books", "read", "Read books"},
		{"books", "write", "Create/Update books"},
		{"books", "delete", "Delete books"},
		{"borrowings", "read", "Read borrowings"},
		{"borrowings", "write", "Create/Update borrowings"},
		{"users", "read", "Read users"},
		{"users", "write", "Create/Update users"},
		{"roles", "read", "Read roles"},
		{"roles", "write", "Create/Update roles"},
		{"permissions", "read", "Read permissions"},
		{"permissions", "write", "Create/Update permissions"},
	}

	for _, perm := range defaultPermissions {
		if _, err := s.permRepo.GetByResourceAndAction(perm.Resource, perm.Action); err != nil {
			// Permission doesn't exist, create it
			newPerm := &domain.Permission{
				Resource:    perm.Resource,
				Action:      perm.Action,
				Description: perm.Description,
			}
			if err := s.permRepo.Create(newPerm); err != nil {
				return fmt.Errorf("failed to create permission %s:%s: %v", perm.Resource, perm.Action, err)
			}
		}
	}

	// Create default roles if they don't exist
	adminRole, err := s.roleRepo.GetByName("admin")
	if err != nil {
		// Create admin role
		adminRole = &domain.Role{
			Name:        "admin",
			Description: "Administrator with full access",
		}
		if err := s.roleRepo.Create(adminRole); err != nil {
			return fmt.Errorf("failed to create admin role: %v", err)
		}
	}

	// Always reassign permissions to admin to ensure it has all permissions
	allPerms, err := s.permRepo.GetAll()
	if err != nil {
		return err
	}
	var adminPermIDs []uint
	for _, perm := range allPerms {
		adminPermIDs = append(adminPermIDs, perm.ID)
	}

	// AssignPermissions already clears existing permissions first
	if err := s.roleRepo.AssignPermissions(adminRole.ID, adminPermIDs); err != nil {
		return fmt.Errorf("failed to assign permissions to admin role: %v", err)
	}

	userRole, err := s.roleRepo.GetByName("user")
	if err != nil {
		// Create user role
		userRole = &domain.Role{
			Name:        "user",
			Description: "Regular user with limited access",
		}
		if err := s.roleRepo.Create(userRole); err != nil {
			return fmt.Errorf("failed to create user role: %v", err)
		}
	}

	// Always reassign basic permissions to user
	basicPermissions := []struct {
		Resource string
		Action   string
	}{
		{"books", "read"},
		{"borrowings", "read"},
		{"borrowings", "write"},
	}

	var userPermIDs []uint
	for _, perm := range basicPermissions {
		permission, err := s.permRepo.GetByResourceAndAction(perm.Resource, perm.Action)
		if err != nil {
			return fmt.Errorf("failed to find permission %s:%s: %v", perm.Resource, perm.Action, err)
		}
		userPermIDs = append(userPermIDs, permission.ID)
	}

	// AssignPermissions already clears existing permissions first
	if err := s.roleRepo.AssignPermissions(userRole.ID, userPermIDs); err != nil {
		return fmt.Errorf("failed to assign permissions to user role: %v", err)
	}

	// Clear and sync role permissions to Casbin
	s.enforcer.ClearPolicy()
	if err := s.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to clear casbin policies: %v", err)
	}

	return s.SyncRolePermissions()
}

// CleanupRBACData removes all RBAC data (useful for testing or reset)
func (s *RBACService) CleanupRBACData() error {
	// Clear all policies from Casbin
	s.enforcer.ClearPolicy()
	if err := s.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to clear policies: %v", err)
	}

	return nil
}

// ResetToDefaults clears everything and sets up default policies again
func (s *RBACService) ResetToDefaults() error {
	// Clear policies
	if err := s.CleanupRBACData(); err != nil {
		return err
	}

	// Setup defaults again
	return s.SetupDefaultPolicies()
}
