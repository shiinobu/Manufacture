# Manufacture Backend API

A Go-based backend application for managing core manufacturing data such as users, customers, suppliers, products, and purchase transactions.

The project is built as a REST API with JWT-based authentication and MySQL persistence. It is organized into controllers, models, helpers, middleware, configuration, and routing layers.

## Highlights

- RESTful HTTP API built with Go
- JWT authentication using HS256
- Protected API routes through authentication middleware
- User management
- Customer management
- Supplier management
- Product management
- Purchase transaction management
- MySQL database integration
- Environment-based configuration
- Centralized database and response helpers
- Automated tests for JWT authentication helpers
- GitHub Actions CI with `go test` and `go vet`

## Architecture

```text
Client
  |
  v
HTTP Request
  |
  v
Router (gorilla/mux)
  |
  v
Auth Middleware
  |
  v
Controllers
  |
  v
Helpers / Models
  |
  v
MySQL
```

The application keeps HTTP routing, authentication, controllers, configuration, and database access separated so the backend remains straightforward to maintain and extend.

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.24.2 |
| HTTP Router | Gorilla Mux |
| Authentication | JWT (`golang-jwt/jwt/v5`) |
| Database | MySQL |
| MySQL Driver | `go-sql-driver/mysql` |
| Configuration | Environment variables / `.env` |
| CORS | Gorilla Handlers |
| Logging | Application logging helper |
| CI | GitHub Actions |

## Project Structure

```text
.
├── app/
│   ├── config/                 # Application configuration
│   ├── controllers/            # HTTP handlers by business domain
│   │   ├── auth/
│   │   ├── customer/
│   │   ├── product/
│   │   ├── purchase/
│   │   ├── supplier/
│   │   └── user/
│   ├── helpers/                # Database, JWT, logging, and response helpers
│   ├── middleware/             # HTTP authentication middleware
│   ├── models/                 # Domain/data models
│   └── routes/                 # API route registration
├── .env.example                # Environment variable template
├── .gitignore
├── .github/workflows/ci.yml    # Automated test and vet pipeline
├── go.mod
├── go.sum
└── main.go
```

## Requirements

Before running the application locally, make sure you have:

- Go 1.24.2 or compatible Go 1.24.x release
- MySQL 8.x or compatible MySQL server
- Git

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/shiinobu/Manufacture.git
cd Manufacture
```

### 2. Configure environment variables

Copy the example environment file:

```bash
cp .env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Update `.env` with your local MySQL credentials and a strong JWT secret.

Example:

```env
DB_HOST=127.0.0.1
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PASS=your_database_password
DB_PORT=3306
SERVER_PORT=8080
JWT_KEY=replace_with_a_long_random_secret
```

> `.env` is intended for local/runtime configuration and is ignored by Git. Do not commit real credentials or secrets.

### 3. Prepare the database

Create the MySQL database configured in `DB_NAME`, then make sure the configured database user has the required permissions.

The application expects the database schema used by the models and controllers in `app/models` and `app/controllers`.

### 4. Download dependencies

```bash
go mod download
```

### 5. Run the application

```bash
go run .
```

The server listens on the port configured by `SERVER_PORT`.

## Authentication

The API uses JWT bearer authentication for protected routes.

Send the token using the `token` HTTP header:

```http
Token: Bearer <jwt-token>
```

The authentication middleware validates the token signature, expiration, claims, and signing method before allowing the request to continue.

JWT signing uses **HS256**, and the signing secret is loaded from the `JWT_KEY` environment variable.

## API Domains

The backend is organized around the following business domains:

| Domain | Purpose |
|---|---|
| Auth | User authentication and token generation |
| User | User management |
| Customer | Customer management |
| Supplier | Supplier management |
| Product | Product management |
| Purchase | Purchase transaction management |

Route definitions are centralized in `app/routes/routes.go`.

## Development

Run the application directly during development:

```bash
go run .
```

Run all tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

## Testing & CI

The repository includes automated tests for the JWT helper, covering:

- successful token generation and parsing
- invalid signing secret
- expired tokens
- unexpected JWT signing methods

GitHub Actions runs the following checks for pushes and pull requests targeting `master`:

```text
go mod download
go test ./...
go vet ./...
```

## Configuration

The following environment variables are required:

| Variable | Description |
|---|---|
| `DB_HOST` | MySQL host |
| `DB_NAME` | MySQL database name |
| `DB_USER` | MySQL username |
| `DB_PASS` | MySQL password |
| `DB_PORT` | MySQL port |
| `SERVER_PORT` | HTTP server port |
| `JWT_KEY` | Secret used to sign JWTs |

## Security Notes

- Never commit `.env` files containing real credentials.
- Use a strong, randomly generated `JWT_KEY` outside development examples.
- JWT tokens are restricted to the HS256 signing method.
- Authentication middleware rejects malformed, invalid, expired, or incorrectly signed tokens.
- Database initialization returns errors instead of terminating the process from deep inside the helper layer.

## Current Status

This project is maintained as a backend portfolio project demonstrating Go REST API development, JWT authentication, MySQL integration, and basic backend engineering practices.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
