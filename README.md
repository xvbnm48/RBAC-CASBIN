# Library API with Dynamic RBAC using Casbin

API perpustakaan yang dibangun dengan clean architecture dan mengimplementasikan **Dynamic Role-Based Access Control (RBAC)** menggunakan Casbin. Admin dapat mengelola role dan permission secara dinamis melalui database.

## 🚀 Fitur Utama

- **Clean Architecture** - Memisahkan business logic dari infrastructure
- **Dynamic RBAC** - Role dan permission disimpan di database
- **Flexible Permission System** - Admin dapat membuat role dan permission baru
- **JWT Authentication** - Secure token-based authentication
- **Multi-Role Support** - User dapat memiliki multiple roles
- **Real-time Permission Checking** - Permission dicek secara real-time
- **RESTful API** - Standard REST endpoints dengan proper HTTP methods

## 🔐 Role & Permission System

### Default Roles

#### Admin Role
- ✅ **Full Access** ke semua resources
- ✅ Mengelola buku (Create, Read, Update, Delete)
- ✅ Mengelola role dan permission
- ✅ Assign role ke user
- ✅ Melihat semua peminjaman

#### User Role
- ✅ Melihat daftar buku
- ✅ Meminjam buku (jika tersedia)
- ✅ Mengembalikan buku
- ✅ Melihat riwayat peminjaman sendiri
- ❌ Tidak dapat mengelola buku
- ❌ Tidak dapat melihat peminjaman user lain

### Dynamic Role Management
- **Create Custom Roles** - Admin dapat membuat role baru
- **Assign Permissions** - Flexible permission assignment
- **Multi-User Assignment** - Satu role bisa diberikan ke banyak user
- **Database-Driven** - Semua policy tersimpan di database

## 📁 Struktur Project (Updated)

```
library-api/
├── cmd/
│   ├── server/          # Main application
│   │   └── main.go
│   └── seed/           # Database seeder with RBAC setup
│       └── main.go
├── configs/            # Configuration files
│   ├── casbin_model.conf
│   └── casbin_policy.csv (now generated dynamically)
├── docs/              # Documentation
│   ├── api_examples.md
│   └── dynamic_rbac_api.md  # NEW: Dynamic RBAC documentation
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── domain/        # Domain models & DTOs (UPDATED)
│   ├── handler/       # HTTP handlers (ADDED: role & permission handlers)
│   ├── middleware/    # Custom middleware (UPDATED: dynamic RBAC)
│   ├── repository/    # Data access layer (ADDED: role & permission repos)
│   ├── service/       # Business logic (ADDED: role & permission services)
│   └── utils/         # Utilities
├── pkg/               # Public packages
│   ├── auth/          # JWT authentication
│   ├── database/      # Database connection (UPDATED: new migrations)
│   └── rbac/          # Dynamic Casbin RBAC (COMPLETELY REFACTORED)
├── scripts/           # Shell scripts
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

## 🆕 New API Endpoints

### Role Management (Admin Only)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/roles` | Get all roles |
| GET | `/api/roles/:id` | Get role by ID |
| POST | `/api/roles` | Create new role |
| PUT | `/api/roles/:id` | Update role |
| DELETE | `/api/roles/:id` | Delete role |
| POST | `/api/roles/assign` | Assign roles to user |
| GET | `/api/users/:id/roles` | Get user's roles |

### Permission Management (Admin Only)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/permissions` | Get all permissions |
| GET | `/api/permissions/:id` | Get permission by ID |
| POST | `/api/permissions` | Create new permission |
| DELETE | `/api/permissions/:id` | Delete permission |

## 🛠 Quick Start

### 1. Prerequisites
- Go 1.21+
- PostgreSQL 12+
- Make (optional)

### 2. Setup Database
```bash
# Create database
createdb library_db

# Or using Docker
docker-compose up postgres -d
```

### 3. Setup Project
```bash
# Clone repository
cd golang-casbin

# Copy environment file
cp .env.example .env

# Install dependencies and setup
make setup

# Or manually:
go mod tidy
go run cmd/seed/main.go  # This now sets up dynamic RBAC
```

### 4. Run Application
```bash
# Using Make
make run

# Or directly
go run cmd/server/main.go

# Using Docker
docker-compose up
```

## 🔑 Default Accounts

Setelah menjalankan seeder:

**Admin Account:**
- Username: `admin`
- Password: `admin123`
- Roles: `admin` (with full permissions)

**User Account:**
- Username: `user`
- Password: `user123`
- Roles: `user` (with limited permissions)

## 📝 Example Usage

### 1. Login as Admin
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### 2. Create Custom Role (e.g., Librarian)
```bash
curl -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "librarian",
    "description": "Library staff with book management access",
    "permissions": [1, 2, 3, 4, 5]
  }'
```

### 3. Assign Role to User
```bash
curl -X POST http://localhost:8080/api/roles/assign \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 2,
    "roles": [3]
  }'
```

### 4. Check User's Roles
```bash
curl -X GET http://localhost:8080/api/users/2/roles \
  -H "Authorization: Bearer <admin_token>"
```

## 📊 Default Permissions Matrix

| Resource | Action | Admin | User | Custom Role |
|----------|--------|--------|------|------------|
| books | read | ✅ | ✅ | ⚙️ Configurable |
| books | write | ✅ | ❌ | ⚙️ Configurable |
| books | delete | ✅ | ❌ | ⚙️ Configurable |
| borrowings | read | ✅ (all) | ✅ (own) | ⚙️ Configurable |
| borrowings | write | ✅ | ✅ | ⚙️ Configurable |
| roles | read | ✅ | ❌ | ⚙️ Configurable |
| roles | write | ✅ | ❌ | ⚙️ Configurable |
| permissions | read | ✅ | ❌ | ⚙️ Configurable |
| permissions | write | ✅ | ❌ | ⚙️ Configurable |

## 🐳 Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up --build

# Scale application
docker-compose up --scale app=3
```

## 📚 Documentation

- **API Documentation**: `docs/dynamic_rbac_api.md`
- **Example Requests**: `docs/api_examples.md`
- **Architecture**: Clean Architecture with RBAC

## 🔧 Development

```bash
# Run with auto-reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Run tests
make test

# Build binary
make build

# Clean build artifacts
make clean
```

## 🎯 Key Improvements in This Version

1. **Dynamic Role Management** - Roles disimpan di database, bukan hardcoded
2. **Flexible Permission System** - Permission dapat dikustomisasi
3. **Multi-Role Support** - User dapat memiliki multiple roles
4. **Real-time Updates** - Changes langsung berlaku tanpa restart
5. **Better Separation of Concerns** - RBAC logic terpisah dari business logic
6. **Enhanced Security** - Permission checking berdasarkan user ID, bukan role token

## 🚀 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Authorization**: Casbin RBAC (Dynamic)
- **Architecture**: Clean Architecture
- **Containerization**: Docker & Docker Compose

## 🤝 Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

---

**🎉 Sekarang sistem RBAC Anda sudah dynamic dan flexible! Admin dapat dengan mudah mengelola role dan permission melalui API.**
