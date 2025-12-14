# Instructions for candidates

This is the Go version of the Payment Gateway challenge. If you haven't already read the [README.md](https://github.com/cko-recruitment/) in the root of this organisation, please do so now. 

## Template structure
```
main.go - a skeleton Payment Gateway API
imposters/ - contains the bank simulator configuration. Don't change this
docs/docs.go - Generated file by Swaggo
.editorconfig - don't change this. It ensures a consistent set of rules for submissions when reformatting code
docker-compose.yml - configures the bank simulator
.goreleaser.yml - Goreleaser configuration
```

Feel free to change the structure of the solution, use a different test library etc.

### Swagger
This template uses Swaggo to autodocument the API and create a Swagger spec. The Swagger UI is available at http://localhost:8090/swagger/index.html .

## Architecture & Design Decisions

### 1. **Layered Architecture**

**Decision**: Implemented a layered architecture with distinct layers: handlers → services → repository/clients.

**Rationale**: 
- **Separation of Concerns**: Each layer has a single, well-defined responsibility
- **Testability**: Layers can be tested independently with mocked dependencies
- **Maintainability**: Changes in one layer don't cascade to others
- **Scalability**: Easy to swap implementations (e.g., replace in-memory storage with a database)

**Structure**:
```
internal/
├── api/           # API setup and routing
├── handlers/      # HTTP request/response handling
├── services/      # Business logic
├── repository/    # Data persistence
├── clients/       # External service communication
├── domain/        # Core business entities
├── dto/           # Data transfer objects
├── validators/    # Input validation logic
└── test/mock/     # Generated mocks for testing
```

### 2. **Dependency Injection**

**Decision**: Manual dependency injection through constructor functions.

**Rationale**:
- **Explicit Dependencies**: Clear visibility of what each component needs
- **No Magic**: No framework overhead, easier to understand and debug
- **Testability**: Simple to inject mocks for unit testing
- **Go Idiomatic**: Follows Go best practices without external DI frameworks

**Implementation**: See `internal/api/dependencies.go:16-27`

### 3. **Domain-Driven Design**

**Decision**: Created a `domain` package with core business entities (`Payment`, `Card`, `PaymentStatus`).

**Rationale**:
- **Business Logic Centralization**: Domain logic lives in domain objects
- **Type Safety**: Strong typing for payment statuses prevents invalid states
- **Encapsulation**: Methods like `GetLastFourDigits()` encapsulate card logic
- **Clear Contracts**: Domain models define the business rules

### 4. **Comprehensive Input Validation**

**Decision**: Dedicated validator package with detailed field-level validation.

**Rationale**:
- **User Experience**: Returns all validation errors at once, not just the first one
- **Security**: Validates card numbers, CVV, expiry dates before processing
- **Business Rules**: Enforces supported currencies (USD, EUR, GBP) and amount constraints
- **Reusability**: Validation logic can be reused across different handlers

**Validation Rules**:
- Card number: 14-19 digits, numeric only
- CVV: 3-4 digits, numeric only
- Expiry date: Must be in the future
- Currency: Must be USD, EUR, or GBP
- Amount: Must be positive

### 5. **Error Handling Strategy**

**Decision**: Implemented custom error types for different failure scenarios.

**Rationale**:
- **Specific Error Handling**: Different errors require different responses
- **Bank Resilience**: Distinguishes between bank unavailability and other errors
- **Status Mapping**: Maps bank errors to appropriate payment statuses
  - Bank unavailable/timeout → Payment status: `Rejected`
  - Bank declined → Payment status: `Declined`
  - Bank authorized → Payment status: `Authorized`

**Custom Errors**:
- `ErrPaymentNotFound`: Payment doesn't exist in repository
- `ErrBankUnavailable`: Bank service is down
- `ErrBankTimeout`: Bank request exceeded timeout
- `ErrBankInvalidRequest`: Invalid request to bank

### 6. **Thread-Safe In-Memory Repository**

**Decision**: Used `sync.RWMutex` for concurrent access to the payment storage.

**Rationale**:
- **Concurrency Safety**: Multiple goroutines can safely access the repository
- **Performance**: Read locks allow concurrent reads, write locks ensure data integrity
- **Simplicity**: In-memory storage is sufficient for the challenge scope
- **Production Ready Pattern**: Easy to replace with a database implementation

**Implementation**: See `/internal/repository/payments_repository.go:15-41`

## Others

### 1. UUID for Payment ID Generation
- Uniqueness: Globally unique across distributed systems.
- Security: Unpredictable, prevents ID enumeration attacks
- Project: Since there was no database, managing the payment ID is easier to implement using an already established library.

### 2. Why CVV is NOT Stored
- Security: CVV only needed for authorization, not future reference
- Liability: Storing CVV increases security risks and regulatory liability

### 3. Why `dto` (Data Transfer Objects) instead of `models`
- Clear Intent: DTOs explicitly represent data structures for API communication (request/response)
- Semantic Accuracy: `Models` is ambiguous and could refer to domain models, view models, or database models

### 4. Mock Generation Strategy
Time Saving: No manual mock implementation needed
Standard Tool: `mockgen` is the official Go mock generation tool

## Project Structure
```
.
├── main.go                          # Application entry point
├── docker-compose.yml               # Bank simulator configuration
├── go.mod                           # Go module dependencies
├── docs/                            # Swagger documentation (auto-generated)
├── imposters/                       # Bank simulator configuration
└── internal/
    ├── api/                         # API setup and dependency injection
    │   ├── api.go                   # Router setup and server lifecycle
    │   └── dependencies.go          # Dependency wiring
    ├── handlers/                    # HTTP handlers
    │   ├── payments_handler.go      # Payment endpoints
    │   ├── ping_handler.go          # Health check
    │   └── swagger_handler.go       # Swagger UI
    ├── services/                    # Business logic
    │   └── payments_service.go      # Payment processing logic
    ├── repository/                  # Data persistence
    │   └── payments_repository.go   # In-memory payment storage
    ├── clients/                     # External service clients
    │   └── bank_client.go           # Bank API client
    ├── domain/                      # Core business entities
    │   └── payment.go               # Payment and Card entities
    ├── dto/                         # Data transfer objects
    │   ├── payment.go               # Payment request/response
    │   ├── bank.go                  # Bank API DTOs
    │   ├── error.go                 # Error response structures
    │   └── pong.go                  # Health check response
    ├── validators/                  # Input validation
    │   └── payments_validator.go    # Payment validation rules
    └── test/mock/                   # Generated mocks
        ├── mock_bank_client.go
        ├── mock_payments_repository.go
        └── mock_payments_service.go
```


## API Usage Examples

### Process Payment
```
curl --location 'localhost:8090/api/payments' \
--header 'Content-Type: application/json' \
--data '{
    "card_number": "5275909408908040",
    "expiry_month": 1,
    "expiry_year": 2026,
    "currency": "USD",
    "amount": 10,
    "cvv": "084"
}'
```

### Retrieve Payment
```
curl --location 'localhost:8090/api/payments/c0b0110d-2cb6-45db-92f8-bca15cab45de' \
--header 'Content-Type: application/json'
```

## Future Enhancements
- Database integration (PostgreSQL/MySQL) to replace in-memory storage
- Authentication and authorization (JWT tokens, API keys)
- Distributed tracing (OpenTelemetry)
- Metrics and monitoring (datadog)
- Structured logs
- Configuration management (environment variables, config files)
- Idempotency keys for duplicate request handling
- Webhook notifications for payment status updates
