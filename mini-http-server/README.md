# mini-http-server

> This is a mini HTTP project which is developed in Golang.

## Features

### 1. Basic HTTP Server

- Built using Go's standard `net/http` package
- Runs on port 8000
- Supports multiple routes and HTTP methods

### 2. Routing

- `/` - Home page
- `/about` - About page
- `/api/users` - RESTful API endpoint for user operations
  - `GET`: Retrieve list of users
  - `POST`: Create new user

### 3. Middleware Support

- Logging middleware that captures:
  - Request start time
  - HTTP method
  - Request path
  - Request duration
  - Request completion status

### 4. Error Handling

- Custom error handler for different HTTP status codes:
  - 404 Not Found
  - 500 Internal Server Error
  - Method Not Allowed
- Structured error logging

### 5. JSON Processing

- JSON request parsing and response formatting
- Supports:
  - Request body parsing
  - JSON serialization
  - JSON deserialization
- Content-Type header management

### 6. Data Structures

- User model with JSON tags:

  ```go
  type User struct {
      ID       int    `json:"id"`
      Name     string `json:"name"`
      Email    string `json:"email"`
      CreateAt string `json:"create_at,omitempty"`
  }
  ```

## API Endpoints

### GET /api/users

Returns a list of users

Response example:

```json
[
  {
    "id": 1,
    "name": "John"
  },
  {
    "id": 2,
    "name": "Jane"
  }
]
```

### POST /api/users

Creates a new user

Request body example:

```json
{
  "id": 1,
  "name": "John",
  "email": "john@example.com"
}
```

Response example:

```json
{
  "id": 1,
  "name": "John",
  "email": "john@example.com",
  "create_at": "2024-03-21T10:30:00Z"
}
```

## Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8000`

## Logging

The server includes detailed logging with:

- Timestamp
- File name and line number
- HTTP method
- Request path
- Request duration
- Error information (when applicable)
