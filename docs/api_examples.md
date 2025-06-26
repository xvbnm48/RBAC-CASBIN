# Library API Postman Collection

## Authentication

### Register User
```bash
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "username": "user1",
  "email": "user1@example.com",
  "password": "password123",
  "role": "user"
}
```

### Register Admin
```bash
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "username": "admin1",
  "email": "admin1@example.com",
  "password": "password123",
  "role": "admin"
}
```

### Login
```bash
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "username": "admin1",
  "password": "password123"
}
```

## Books (Admin Only)

### Create Book
```bash
POST http://localhost:8080/api/books
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "title": "The Go Programming Language",
  "author": "Alan A. A. Donovan, Brian W. Kernighan",
  "isbn": "978-0134190440",
  "description": "A comprehensive guide to Go programming",
  "stock": 5
}
```

### Get All Books
```bash
GET http://localhost:8080/api/books
Authorization: Bearer <your-jwt-token>
```

### Get Book by ID
```bash
GET http://localhost:8080/api/books/1
Authorization: Bearer <your-jwt-token>
```

### Update Book
```bash
PUT http://localhost:8080/api/books/1
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "title": "Updated Title",
  "stock": 10
}
```

### Delete Book
```bash
DELETE http://localhost:8080/api/books/1
Authorization: Bearer <your-jwt-token>
```

## Borrowing

### Borrow Book
```bash
POST http://localhost:8080/api/borrowings
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "book_id": 1
}
```

### Get Borrowings
```bash
GET http://localhost:8080/api/borrowings
Authorization: Bearer <your-jwt-token>
```

### Return Book
```bash
PUT http://localhost:8080/api/borrowings/1/return
Authorization: Bearer <your-jwt-token>
```
