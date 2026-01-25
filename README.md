# Customer Onboarding Service

A RESTful microservice API built with **Go** and the **Gin Framework** for customer onboarding with phone verification via OTP. This service includes Swagger documentation, MS SQL Server integration with GORM, and comprehensive unit tests.

## 📋 Table of Contents

- [Features](#features)
- [Requirements Met](#requirements-met)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running Tests](#running-tests)
- [Swagger Documentation](#swagger-documentation)
- [Architecture & Design Patterns](#architecture--design-patterns)

## ✨ Features

- **Customer Onboarding** with phone, email, password, state, and LGA
- **Phone Verification via OTP** before onboarding is complete
- **State-LGA Validation** ensuring LGA belongs to selected state
- **Mock OTP Sending** (ready for SMS provider integration)
- **Get All Customers** endpoint
- **Get Banks** from external API integration
- **Swagger UI** documentation
- **MS SQL Server** database with GORM ORM
- **In-Memory Fallback** when database is not configured
- **Unit Tests** with mocking for all layers

## ✅ Requirements Met

| Requirement | Status | Implementation |
|-------------|--------|----------------|
| Onboard customer endpoint | ✅ | `POST /api/v1/customers/onboard` |
| Phone verification via OTP | ✅ | `POST /api/v1/customers/verify-otp` |
| LGA mapped to state | ✅ | Validated in `IsValidLGA()` |
| Mock OTP sending | ✅ | `MockSendOTP()` logs to console |
| Get all customers | ✅ | `GET /api/v1/customers` |
| Get banks endpoint | ✅ | `GET /api/v1/banks` |
| Go + Gin Framework | ✅ | Using Gin v1.11.0 |
| Swagger documentation | ✅ | Available at `/swagger/index.html` |
| MS SQL DB + ORM | ✅ | GORM with SQL Server driver |
| Unit testing | ✅ | 25+ tests with testify |
| Engineering principles | ✅ | Clean architecture, DI, interfaces |

## 📁 Project Structure

```
gin-quickstart/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point with Swagger annotations
├── docs/
│   ├── docs.go                     # Generated Swagger docs
│   ├── swagger.json                # OpenAPI spec (JSON)
│   └── swagger.yaml                # OpenAPI spec (YAML)
├── internal/
│   ├── clients/
│   │   └── bank_client.go          # External bank API client
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── controller/
│   │   └── customer_controller.go  # HTTP handlers with Swagger annotations
│   ├── database/
│   │   └── database.go             # GORM database connection
│   ├── models/
│   │   ├── customer.go             # Customer & OTPRecord entities
│   │   └── dto.go                  # Request/Response DTOs
│   ├── repositories/
│   │   ├── interfaces.go           # Repository interfaces
│   │   ├── customer_repository.go  # Customer data access
│   │   ├── otp_repository.go       # OTP data access
│   │   └── repositories_test.go    # Repository tests
│   ├── routes/
│   │   └── routes.go               # Route definitions
│   ├── services/
│   │   ├── customer_service.go     # Business logic
│   │   └── customer_service_test.go # Service tests
│   └── utils/
│       ├── lga_validator.go        # State/LGA validation
│       ├── otp.go                  # OTP generation utilities
│       └── utils_test.go           # Utility tests
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- MS SQL Server (optional - falls back to in-memory storage)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd gin-quickstart
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Generate Swagger docs (optional - already generated):**
   ```bash
   swag init -g cmd/server/main.go -o docs
   ```

4. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

The server starts on `http://localhost:8080` by default.

## 📡 API Endpoints

### Health Check

```http
GET /ping
```

**Response:**
```json
{
  "success": true,
  "message": "pong"
}
```

---

### 1. Onboard Customer

```http
POST /api/v1/customers/onboard
```

**Request Body:**
```json
{
  "phone": "+2348012345678",
  "email": "customer@example.com",
  "password": "securepassword123",
  "state": "Lagos",
  "lga": "Ikeja"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "Customer registered. OTP sent to phone for verification."
}
```

**Error Responses:**
- `400` - Invalid LGA for state, phone/email already exists
- `400` - Invalid request payload

---

### 2. Verify OTP (Complete Onboarding)

```http
POST /api/v1/customers/verify-otp
```

**Request Body:**
```json
{
  "phone": "+2348012345678",
  "otp": "123456"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Phone number verified. Onboarding complete!"
}
```

**Error Responses:**
- `400` - Invalid or expired OTP
- `400` - Phone not registered

---

### 3. Resend OTP

```http
POST /api/v1/customers/resend-otp
```

**Request Body:**
```json
{
  "phone": "+2348012345678"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "OTP resent successfully"
}
```

---

### 4. Get All Customers

```http
GET /api/v1/customers
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Customers retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "phone": "+2348012345678",
      "email": "customer@example.com",
      "state": "Lagos",
      "lga": "Ikeja",
      "isPhoneVerified": true,
      "createdAt": "2026-01-25T22:00:00Z"
    }
  ]
}
```

---

### 5. Get Banks

```http
GET /api/v1/banks
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Banks retrieved successfully",
  "data": [
    {
      "code": "044",
      "name": "Access Bank"
    }
  ]
}
```

## ⚙️ Configuration

Configure via environment variables:

| Variable      | Default            | Description                     |
|---------------|--------------------|---------------------------------|
| `SERVER_PORT` | `8080`             | Port the server listens on      |
| `GIN_MODE`    | `debug`            | Gin mode (debug/release/test)   |
| `DB_HOST`     | `localhost`        | MS SQL Server host              |
| `DB_PORT`     | `1433`             | MS SQL Server port              |
| `DB_USER`     | `sa`               | Database username               |
| `DB_PASSWORD` | *(empty)*          | Database password               |
| `DB_NAME`     | `customer_onboarding` | Database name                |

**Example:**
```bash
export DB_HOST=localhost
export DB_PORT=1433
export DB_USER=sa
export DB_PASSWORD=YourStrong@Password
export DB_NAME=customer_onboarding
export GIN_MODE=release
go run cmd/server/main.go
```

## 🗄️ Database Setup

### MS SQL Server

1. Create the database:
   ```sql
   CREATE DATABASE customer_onboarding;
   ```

2. The application will auto-migrate tables on startup.

### In-Memory Mode (Development)

If `DB_PASSWORD` is not set, the application falls back to in-memory storage. This is useful for development and testing but data is lost on restart.

## 🧪 Running Tests

Run all tests:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -v -cover
```

Run specific package tests:
```bash
go test ./internal/services/... -v
go test ./internal/repositories/... -v
go test ./internal/utils/... -v
```

### Test Coverage

| Package | Tests | Coverage |
|---------|-------|----------|
| services | 10 | ~80% |
| repositories | 10 | ~80% |
| utils | 12 | ~80% |

## 📚 Swagger Documentation

Swagger UI is available at:
```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger docs after API changes:
```bash
swag init -g cmd/server/main.go -o docs
```

## 🏗️ Architecture & Design Patterns

### Clean Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Request                         │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                   Routes (DI Container)                  │
│         Dependency injection & route registration        │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                      Controller                          │
│    HTTP handling, validation, request/response mapping   │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                       Service                            │
│           Business logic, orchestration, validation      │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                     Repository                           │
│    Data access abstraction (MS SQL / In-Memory)          │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────┐
│                      Database                            │
│                   MS SQL Server                          │
└─────────────────────────────────────────────────────────┘
```

### Design Patterns Applied

1. **Repository Pattern** - Abstracts data access with interfaces
2. **Dependency Injection** - Services receive dependencies through constructors
3. **Interface Segregation** - Small, focused interfaces for repositories
4. **DTO Pattern** - Separate models for API requests/responses
5. **Factory Pattern** - Constructor functions for creating instances

### SOLID Principles

- **S**ingle Responsibility - Each layer has one job
- **O**pen/Closed - Extend via interfaces without modifying
- **L**iskov Substitution - Repositories are interchangeable
- **I**nterface Segregation - Small, focused interfaces
- **D**ependency Inversion - Depend on abstractions (interfaces)

## 🔒 Supported States & LGAs

| State  | Local Government Areas      |
|--------|----------------------------|
| Lagos  | Ikeja, Epe, Ikorodu        |
| Oyo    | Ibadan North, Ibadan South |

## 📝 License

This project is for demonstration purposes.

---

## 🛠️ Future Improvements

- [ ] Integrate real SMS provider (Twilio, Termii, etc.)
- [ ] Add rate limiting for OTP endpoints
- [ ] Implement JWT authentication
- [ ] Add more Nigerian states and LGAs
- [ ] Add API versioning headers
- [ ] Implement request logging middleware
- [ ] Add containerization with Docker
