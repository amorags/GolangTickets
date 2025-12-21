# Ticket Booking API

A RESTful API built with Go for managing event tickets, featuring JWT-based authentication and PostgreSQL database integration.

## Features

- User authentication (signup/login) with JWT tokens
- Password hashing with bcrypt
- PostgreSQL database with GORM ORM
- RESTful API design with Chi router
- Docker and Docker Compose support
- Protected routes with middleware
- Event and ticket management models

## Tech Stack

- **Go** 1.25.5
- **Chi** v5 - HTTP router
- **GORM** - ORM library
- **PostgreSQL** 16 - Database
- **JWT** - Authentication
- **Docker** - Containerization

## Project Structure

```
golang_test/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handlers/
│   │   └── auth_handler.go      # Authentication handlers
│   ├── middleware/
│   │   └── auth_middleware.go   # JWT authentication middleware
│   ├── models/
│   │   └── models.go            # Database models (User, Event, Ticket)
│   ├── repository/
│   │   ├── database.go          # Database connection
│   │   └── user_repository.go   # User data access
│   └── utils/
│       ├── jwt.go               # JWT token utilities
│       ├── password.go          # Password hashing utilities
│       └── response.go          # HTTP response helpers
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker build configuration
├── go.mod                       # Go module dependencies
└── go.sum                       # Dependency checksums
```

## Prerequisites

- Go 1.25.5 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL 16 (if running without Docker)

## Setup

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/alexs/golang_test.git
cd golang_test
```

2. Start the application with Docker Compose:
```bash
docker-compose up --build
```

This will:
- Start a PostgreSQL database on port 5432
- Build and run the API on port 8080
- Automatically configure the database connection

### Manual Setup

1. Clone the repository:
```bash
git clone https://github.com/alexs/golang_test.git
cd golang_test
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database and create a database named `ticket_db`

4. Set environment variables:
```bash
export DB_HOST=localhost
export DB_USER=user
export DB_PASSWORD=password
export DB_NAME=ticket_db
export DB_PORT=5432
export JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-min-32-chars
```

5. Run the application:
```bash
go run cmd/api/main.go
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_USER` | Database user | `user` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `ticket_db` |
| `DB_PORT` | Database port | `5432` |
| `JWT_SECRET` | JWT signing secret (min 32 chars) | *Required* |

## API Endpoints

### Public Endpoints

#### Health Check
```
GET /health
```
Returns the API health status.

#### User Signup
```
POST /auth/signup
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### User Login
```
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```
Returns a JWT token for authentication.

### Protected Endpoints

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

#### Get User Profile
```
GET /profile
```
Returns the authenticated user's profile information.

## Database Models

### User
- `ID` - Primary key
- `Username` - Unique username
- `Email` - Unique email address
- `Password` - Bcrypt hashed password
- `CreatedAt`, `UpdatedAt`, `DeletedAt` - GORM timestamps

### Event
- `ID` - Primary key
- `Name` - Event name
- `Description` - Event description
- `Location` - Event location
- `Date` - Event date and time
- `Tickets` - Related tickets (one-to-many)
- `CreatedAt`, `UpdatedAt`, `DeletedAt` - GORM timestamps

### Ticket
- `ID` - Primary key
- `EventID` - Foreign key to Event
- `Price` - Ticket price
- `Seat` - Seat identifier
- `IsSold` - Sale status (default: false)
- `CreatedAt`, `UpdatedAt`, `DeletedAt` - GORM timestamps

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/api cmd/api/main.go
```

### Docker Build
```bash
docker build -t golang-ticket-api .
docker run -p 8080:8080 --env-file .env golang-ticket-api
```

## Security Considerations

- Passwords are hashed using bcrypt before storage
- JWT tokens expire after 24 hours
- Change the `JWT_SECRET` in production to a strong, random value (minimum 32 characters)
- Database credentials should be stored securely (use secrets management in production)

## Future Enhancements

- Event management endpoints (CRUD operations)
- Ticket booking and purchase functionality
- Payment integration
- User ticket history
- Admin role and permissions
- Rate limiting
- API documentation with Swagger/OpenAPI

## License

This project is provided as-is for educational and development purposes.
