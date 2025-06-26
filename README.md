# Library API with Clean Architecture and Casbin RBAC

API perpustakaan yang dibangun dengan clean architecture dan mengimplementasikan Casbin untuk Role-Based Access Control (RBAC).

## Fitur

- **Clean Architecture** - Memisahkan business logic dari infrastructure
- **Role-Based Access Control (RBAC)** - Menggunakan Casbin untuk authorization
- **JWT Authentication** - Secure token-based authentication
- **User Management** - Admin dan User dengan permission berbeda
- **Book Management** - CRUD operations untuk buku
- **Borrowing System** - Sistem peminjaman dengan tracking
- **RESTful API** - Standard REST endpoints
- **Database Migration** - Auto-migration dengan GORM
- **Docker Support** - Containerization ready

## Roles & Permissions

### Admin
- ✅ Mengelola buku (Create, Read, Update, Delete)
- ✅ Melihat semua peminjaman
- ✅ Mengelola user
- ✅ Akses ke semua endpoint

### User
- ✅ Melihat daftar buku
- ✅ Meminjam buku (jika tersedia)
- ✅ Mengembalikan buku
- ✅ Melihat riwayat peminjaman sendiri
- ❌ Tidak dapat mengelola buku
- ❌ Tidak dapat melihat peminjaman user lain

## Struktur Project

```
library-api/
├── cmd/
│   ├── server/          # Main application
│   │   └── main.go
│   └── seed/           # Database seeder
│       └── main.go
├── configs/            # Configuration files
│   ├── casbin_model.conf
│   └── casbin_policy.csv
├── docs/              # Documentation
│   └── api_examples.md
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── domain/        # Domain models & DTOs
│   ├── handler/       # HTTP handlers
│   ├── middleware/    # Custom middleware
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic
│   └── utils/         # Utilities
├── pkg/               # Public packages
│   ├── auth/          # JWT authentication
│   ├── database/      # Database connection
│   └── rbac/          # Casbin RBAC
├── scripts/           # Shell scripts
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

## Quick Start

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
git clone <your-repo>
cd golang-casbin

# Copy environment file
cp .env.example .env

# Install dependencies and setup
make setup

# Or manually:
go mod tidy
go run cmd/seed/main.go
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

## Environment Variables

```bash
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=library_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login user |

### Books
| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/books` | Get all books | user, admin |
| GET | `/api/books/:id` | Get book by ID | user, admin |
| POST | `/api/books` | Create new book | admin only |
| PUT | `/api/books/:id` | Update book | admin only |
| DELETE | `/api/books/:id` | Delete book | admin only |

### Borrowing
| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/borrowings` | Get borrowings | user (own), admin (all) |
| POST | `/api/borrowings` | Borrow book | user, admin |
| PUT | `/api/borrowings/:id/return` | Return book | user, admin |

## Default Accounts

Setelah menjalankan seeder:

**Admin Account:**
- Username: `admin`
- Password: `admin123`
- Role: `admin`

**User Account:**
- Username: `user`
- Password: `user123`
- Role: `user`

## Testing API

Lihat file `docs/api_examples.md` untuk contoh request/response atau gunakan collection Postman yang tersedia.

## Development

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

## Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up --build

# Run only database
docker-compose up postgres

# Scale application
docker-compose up --scale app=3
```

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Authorization**: Casbin RBAC
- **Containerization**: Docker & Docker Compose
