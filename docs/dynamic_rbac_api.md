# API Documentation - Dynamic RBAC Library System

## Overview

This API provides role-based access control (RBAC) functionality where administrators can dynamically assign roles and permissions to users. The system stores roles and permissions in the database, making it flexible and manageable.

## Authentication

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Default Roles

### Admin Role
- Full access to all resources
- Can manage books (CRUD)
- Can manage roles and permissions
- Can assign roles to users
- Can view all borrowings

### User Role
- Can view books
- Can borrow and return books
- Can view own borrowings only

## API Endpoints

### Authentication

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@library.com",
    "role": "admin"
  }
}
```

### Role Management (Admin Only)

#### Create Role
```http
POST /api/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "librarian",
  "description": "Library staff with book management access",
  "permissions": [1, 2, 3]
}
```

#### Get All Roles
```http
GET /api/roles
Authorization: Bearer <token>
```

#### Get Role by ID
```http
GET /api/roles/:id
Authorization: Bearer <token>
```

#### Update Role
```http
PUT /api/roles/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "updated_librarian",
  "description": "Updated description",
  "permissions": [1, 2, 3, 4]
}
```

#### Delete Role
```http
DELETE /api/roles/:id
Authorization: Bearer <token>
```

#### Assign Roles to User
```http
POST /api/roles/assign
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 2,
  "roles": [1, 2]
}
```

#### Get User Roles
```http
GET /api/users/:id/roles
Authorization: Bearer <token>
```

Response:
```json
{
  "id": 2,
  "username": "john_doe",
  "email": "john@example.com",
  "roles": [
    {
      "id": 1,
      "name": "user",
      "description": "Regular user with limited access",
      "permissions": [
        {
          "id": 1,
          "resource": "books",
          "action": "read",
          "description": "Read books"
        }
      ]
    }
  ]
}
```

### Permission Management (Admin Only)

#### Create Permission
```http
POST /api/permissions
Authorization: Bearer <token>
Content-Type: application/json

{
  "resource": "reports",
  "action": "read",
  "description": "Read reports"
}
```

#### Get All Permissions
```http
GET /api/permissions
Authorization: Bearer <token>
```

#### Get Permission by ID
```http
GET /api/permissions/:id
Authorization: Bearer <token>
```

#### Delete Permission
```http
DELETE /api/permissions/:id
Authorization: Bearer <token>
```

### Book Management

#### Get All Books
```http
GET /api/books
Authorization: Bearer <token>
```

#### Get Book by ID
```http
GET /api/books/:id
Authorization: Bearer <token>
```

#### Create Book (Admin Only)
```http
POST /api/books
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Clean Architecture",
  "author": "Robert C. Martin",
  "isbn": "978-0134494166",
  "description": "A guide to software architecture",
  "stock": 10
}
```

#### Update Book (Admin Only)
```http
PUT /api/books/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Clean Architecture",
  "stock": 15
}
```

#### Delete Book (Admin Only)
```http
DELETE /api/books/:id
Authorization: Bearer <token>
```

### Borrowing System

#### Get Borrowings
```http
GET /api/borrowings
Authorization: Bearer <token>
```
- Admin: Returns all borrowings
- User: Returns only own borrowings

#### Borrow Book
```http
POST /api/borrowings
Authorization: Bearer <token>
Content-Type: application/json

{
  "book_id": 1
}
```

#### Return Book
```http
PUT /api/borrowings/:id/return
Authorization: Bearer <token>
```

## Default Permissions

| Resource | Action | Description |
|----------|--------|-------------|
| books | read | View books |
| books | write | Create/Update books |
| books | delete | Delete books |
| borrowings | read | View borrowings |
| borrowings | write | Create/Update borrowings |
| users | read | View users |
| users | write | Create/Update users |
| roles | read | View roles |
| roles | write | Create/Update roles |
| permissions | read | View permissions |
| permissions | write | Create/Update permissions |

## Setup and Usage

### 1. Initial Setup
```bash
# Run database migration and seed
go run cmd/seed/main.go
```

### 2. Default Accounts
- **Admin**: username: `admin`, password: `admin123`
- **User**: username: `user`, password: `user123`

### 3. Creating Custom Roles
1. Login as admin
2. Create permissions if needed
3. Create role with assigned permissions
4. Assign role to users

### 4. Example: Creating Librarian Role
```bash
# 1. Create role
curl -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "librarian",
    "description": "Library staff",
    "permissions": [1, 2, 3, 4, 5]
  }'

# 2. Assign role to user
curl -X POST http://localhost:8080/api/roles/assign \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 2,
    "roles": [3]
  }'
```

## Error Responses

```json
{
  "error": "Insufficient permissions"
}
```

```json
{
  "error": "Role with this name already exists"
}
```

```json
{
  "error": "User not found"
}
```

## Features

- ✅ Dynamic role creation and management
- ✅ Flexible permission system
- ✅ Database-stored policies
- ✅ Role assignment to users
- ✅ Real-time permission checking
- ✅ Multiple roles per user support
- ✅ Clean separation of concerns
- ✅ RESTful API design
