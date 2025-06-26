# Migration Guide: Static to Dynamic RBAC

## Overview

Project ini telah di-refactor dari sistem RBAC statis (berbasis file) menjadi sistem RBAC dinamis (berbasis database). Berikut adalah perubahan-perubahan yang dilakukan.

## ✨ Perubahan Utama

### 1. Database Schema Changes

**Tabel Baru:**
- `roles` - Menyimpan role yang dapat dibuat secara dinamis
- `permissions` - Menyimpan permission yang fleksibel
- `user_roles` - Mapping many-to-many antara user dan role
- `role_permissions` - Mapping many-to-many antara role dan permission

**Tabel yang Dimodifikasi:**
- `users` - Masih memiliki kolom `role` untuk backward compatibility

### 2. Architecture Changes

**Sebelum (Static RBAC):**
```
configs/casbin_policy.csv (static policies)
     ↓
rbac.SetupDefaultPolicies()
     ↓
Casbin Enforcer checks role directly
```

**Sesudah (Dynamic RBAC):**
```
Database Tables (roles, permissions, user_roles)
     ↓
RBACService.CheckPermission()
     ↓
Dynamic policy generation & checking
```

### 3. New Components

#### Repositories
- `RoleRepository` - CRUD operations untuk roles
- `PermissionRepository` - CRUD operations untuk permissions  
- `UserRoleRepository` - Manage user-role assignments

#### Services
- `RoleService` - Business logic untuk role management
- `PermissionService` - Business logic untuk permission management

#### Handlers
- `RoleHandler` - HTTP endpoints untuk role management
- `PermissionHandler` - HTTP endpoints untuk permission management

#### RBAC Service
- `pkg/rbac/RBACService` - Centralized RBAC logic dengan database integration

## 🔄 Migration Steps

### 1. Database Migration
Sistem akan otomatis membuat tabel baru saat pertama kali dijalankan:
```bash
go run cmd/seed/main.go
```

### 2. Default Setup
Seeder akan secara otomatis:
- Membuat default permissions
- Membuat default roles (admin, user)
- Assign permissions ke roles
- Assign roles ke existing users

### 3. Backward Compatibility
- Kolom `users.role` masih ada untuk kompatibilitas
- Existing users akan otomatis mendapat role assignment

## 🆕 New Capabilities

### 1. Dynamic Role Creation
```bash
curl -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "librarian",
    "description": "Library staff",
    "permissions": [1, 2, 3, 4, 5]
  }'
```

### 2. Flexible Permission Assignment
```bash
curl -X PUT http://localhost:8080/api/roles/3 \
  -H "Authorization: Bearer <token>" \
  -d '{
    "permissions": [1, 2, 3, 4, 5, 6, 7]
  }'
```

### 3. Multi-Role Support
```bash
curl -X POST http://localhost:8080/api/roles/assign \
  -H "Authorization: Bearer <token>" \
  -d '{
    "user_id": 2,
    "roles": [1, 3]  # User bisa punya multiple roles
  }'
```

### 4. Real-time Permission Updates
- Perubahan permission langsung berlaku tanpa restart
- No need to modify code untuk permission baru

## 📋 Migration Checklist

### Pre-Migration
- [x] Backup existing database
- [x] Test new system di development environment
- [x] Verify all existing functionality still works

### Migration Process
- [x] Deploy new code
- [x] Run database migrations (`go run cmd/seed/main.go`)
- [x] Verify default roles and permissions created
- [x] Test authentication and authorization

### Post-Migration
- [x] Test all existing API endpoints
- [x] Test new role management endpoints
- [x] Verify user permissions work correctly
- [x] Update documentation

## 🐛 Troubleshooting

### Issue: User can't access previously accessible endpoints
**Solution:** Check if user has been assigned appropriate roles:
```bash
curl -X GET http://localhost:8080/api/users/{user_id}/roles \
  -H "Authorization: Bearer <admin_token>"
```

### Issue: New permissions not working
**Solution:** Sync permissions to Casbin:
```go
rbacService.SyncRolePermissions()
```

### Issue: Role assignment not working
**Solution:** Verify role and user exist:
```bash
# Check if role exists
curl -X GET http://localhost:8080/api/roles \
  -H "Authorization: Bearer <token>"

# Check if user exists
curl -X GET http://localhost:8080/api/users/{user_id}/roles \
  -H "Authorization: Bearer <token>"
```

## 📈 Benefits

### 1. Flexibility
- No need to modify code untuk permission baru
- Admin dapat membuat role custom sesuai kebutuhan
- Dynamic permission assignment

### 2. Scalability  
- Support untuk multiple roles per user
- Granular permission control
- Database-driven policies

### 3. Maintainability
- Clear separation of concerns
- Centralized RBAC logic
- Easy to test and debug

### 4. Security
- Permission checking berdasarkan user ID
- Real-time permission validation
- Audit trail untuk role changes

## 🎯 Next Steps

1. **Custom Role Creation** - Create roles sesuai kebutuhan bisnis
2. **Permission Auditing** - Implement logging untuk permission changes
3. **Role Templates** - Create template roles untuk use cases umum
4. **API Documentation** - Update API docs dengan new endpoints
5. **Testing** - Add comprehensive tests untuk RBAC functionality

## 📝 Code Examples

### Before (Static)
```go
// Hard-coded policy check
allowed, _ := enforcer.Enforce(role, resource, action)
```

### After (Dynamic)
```go
// Dynamic permission check based on user's roles
allowed, _ := rbacService.CheckPermission(userID, resource, action)
```

### Creating Custom Role
```go
roleService.CreateRole(&domain.CreateRoleRequest{
    Name: "manager",
    Description: "Department manager",
    Permissions: []uint{1, 2, 3, 5, 8},
})
```

---

**🎉 Migration selesai! Sistem RBAC Anda sekarang sudah dynamic dan siap untuk scale!**
