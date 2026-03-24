# Todos Open API

A RESTful API for managing users with in-memory storage. Built with Go using the standard `net/http` package.

## Features

- User CRUD operations (Create, Read, Update, Delete)
- In-memory storage (no database required)
- Structured JSON logging
- Graceful shutdown handling
- Request logging middleware
- Health check endpoints

## Quick Start

### Prerequisites

- Go 1.21 or later

### Running the Server

```bash
# Clone the repository
git clone https://github.com/MrEraba/todos-open-api.git
cd todos-open-api

# Run the server
go run main.go
```

The server will start on `http://localhost:8080` by default.

### Configuration

You can configure the server address using the `ADDR` environment variable:

```bash
ADDR=:3000 go run main.go
```

## API Endpoints

### Health Check Endpoints

#### GET /health
Check if the API is running.

**Response:**
```json
{
  "status": "healthy",
  "service": "todos-api"
}
```

#### GET /status
Get API status including user count.

**Response:**
```json
{
  "status": "running",
  "service": "todos-api",
  "user_count": 5,
  "server_info": "Go HTTP Server"
}
```

### User Endpoints

#### GET /users
List all users.

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

#### GET /users/{id}
Get a user by ID.

**Parameters:**
- `id` (path) - User UUID

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "not_found",
  "message": "User not found"
}
```

#### POST /users
Create a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "active": false,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Response:** `400 Bad Request`
```json
{
  "error": "validation_error",
  "message": "Name is required"
}
```

**Error Response:** `409 Conflict`
```json
{
  "error": "duplicate",
  "message": "A user with this email already exists"
}
```

#### PUT /users/{id}
Update an existing user.

**Parameters:**
- `id` (path) - User UUID

**Request Body:** (all fields optional)
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "active": true
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "not_found",
  "message": "User not found"
}
```

#### DELETE /users/{id}
Delete a user.

**Parameters:**
- `id` (path) - User UUID

**Response:** `204 No Content`

**Error Response:** `404 Not Found`
```json
{
  "error": "not_found",
  "message": "User not found"
}
```

## Example Usage with cURL

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# List all users
curl http://localhost:8080/users

# Get a specific user
curl http://localhost:8080/users/{user-id}

# Update a user
curl -X PUT http://localhost:8080/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{"active": true}'

# Delete a user
curl -X DELETE http://localhost:8080/users/{user-id}

# Health check
curl http://localhost:8080/health
```

## Project Structure

```
.
├── main.go           # Application entry point
├── models/
│   └── users.go      # User model and request/response types
├── handlers/
│   └── user_handler.go  # HTTP handlers for user endpoints
├── store/
│   └── user_store.go    # In-memory user storage
├── go.mod
└── README.md
```

## API Design Principles

1. **RESTful Conventions**: Uses standard HTTP methods and status codes
2. **JSON Responses**: All responses are JSON formatted
3. **Error Handling**: Consistent error response format with error codes
4. **Validation**: Basic request validation with clear error messages
5. **Thread Safety**: In-memory store uses mutex locks for concurrent access
6. **Graceful Shutdown**: Server handles SIGINT/SIGTERM for clean shutdown

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (successful delete) |
| 400 | Bad Request (validation error) |
| 404 | Not Found |
| 409 | Conflict (duplicate resource) |
| 500 | Internal Server Error |

## Limitations

- **In-Memory Storage**: Data is not persisted between server restarts
- **No Authentication**: API has no authentication/authorization
- **Single Instance**: Not designed for horizontal scaling

## Future Enhancements

- [ ] Add authentication (JWT, OAuth)
- [ ] Add pagination for list endpoints
- [ ] Add filtering and sorting options
- [ ] Add request rate limiting
- [ ] Add input sanitization
- [ ] Add database persistence option
- [ ] Add OpenAPI/Swagger documentation
- [ ] Add request validation middleware

## License

MIT
