# Consultancy API

A simple REST API built in Go that manages consultants, skills, and projects. This project demonstrates Go's concurrency patterns using goroutines, channels, and synchronization primitives.

## Features

- RESTful API endpoints for consultants, skills, and projects
- In-memory data store with thread-safe operations
- Concurrent data processing with goroutines and channels
- Graceful server shutdown
- Logging middleware

## Prerequisites

- Go 1.16 or later
- [Gorilla Mux](https://github.com/gorilla/mux) for routing

## Project Structure

## Installation and Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/consultancy-api.git
   cd consultancy-api

go mod init github.com/yourusername/consultancy-api
go mod tidy



## Installation and Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/consultancy-api.git
   cd consultancy-api

Update import paths:

Replace github.com/yourusername/consultancy-api with your actual module path in all files


Initialize the Go module:
bashgo mod init github.com/yourusername/consultancy-api
go mod tidy

Run the server:
bashgo run main.go


The server will start on http://localhost:8080
API Endpoints
Consultants

GET /api/consultants - Get all consultants
GET /api/consultants/{id} - Get a specific consultant
POST /api/consultants - Create a new consultant
PUT /api/consultants/{id} - Update a consultant
DELETE /api/consultants/{id} - Delete a consultant
GET /api/consultants/skills/{skill_id} - Get consultants with a specific skill
GET /api/consultants/projects/{project_id} - Get consultants assigned to a specific project

Skills

GET /api/skills - Get all skills
GET /api/skills/{id} - Get a specific skill
POST /api/skills - Create a new skill
PUT /api/skills/{id} - Update a skill
DELETE /api/skills/{id} - Delete a skill

Projects

GET /api/projects - Get all projects
GET /api/projects/{id} - Get a specific project
POST /api/projects - Create a new project
PUT /api/projects/{id} - Update a project
DELETE /api/projects/{id} - Delete a project
GET /api/projects/{id}/details - Get a project with consultant and skill details

Testing API Endpoints
Using curl
Get all consultants:
bashcurl -X GET http://localhost:8080/api/consultants
Create a new consultant:
bashcurl -X POST http://localhost:8080/api/consultants \
-H "Content-Type: application/json" \
-d '{"name":"Alice Cooper","email":"alice@example.com","skill_ids":[1,2]}'
Get a specific project with details:
bashcurl -X GET http://localhost:8080/api/projects/1/details
Using Postman

Import the following collection:
https://www.getpostman.com/collections/[collection_id]
(Replace with an actual collection URL if you create one)
Test each endpoint in the Postman interface.

Concurrency Features
The API demonstrates several Go concurrency patterns:

Goroutines: Used in the server startup, metrics collection, and the project details endpoint.
Channels: Used for communication between goroutines, especially in the project details handler.
WaitGroups: Used to coordinate multiple goroutines and ensure they all complete before proceeding.
Mutex Locks: Used in the data store to provide thread-safe operations.
Graceful Shutdown: Implements a pattern for gracefully shutting down the server with context.

Extending the Project
To extend this project:

Add Authentication: Implement JWT authentication middleware.
Persistence: Replace the in-memory store with a database like PostgreSQL.
Testing: Add unit and integration tests.
API Documentation: Add Swagger documentation.

License
This project is licensed under the MIT License - see the LICENSE file for details.

## Go Mod File


# Go Programming: Building Production-Ready REST APIs

## Table of Contents

1. [Go Language Fundamentals](#1-go-language-fundamentals)
   - [Types and Data Structures](#types-and-data-structures)
   - [Packages and Imports](#packages-and-imports)
   - [Error Handling Philosophy](#error-handling-philosophy)
   - [Functions and Methods](#functions-and-methods)

2. [Concurrency in Go](#2-concurrency-in-go)
   - [Goroutines](#goroutines)
   - [Channels](#channels)
   - [WaitGroups](#waitgroups)
   - [Mutex](#mutex)
   - [Context Package](#context-package)

3. [Building REST APIs](#3-building-rest-apis)
   - [HTTP Package](#http-package)
   - [Router/Mux](#routermux)
   - [Handlers](#handlers)
   - [Middleware](#middleware)
   - [Request/Response Processing](#requestresponse-processing)

4. [Database Integration](#4-database-integration)
   - [SQL Package](#sql-package)
   - [Connection Pooling](#connection-pooling)
   - [Transactions](#transactions)
   - [Repository Pattern](#repository-pattern)

5. [Project Organization](#5-project-organization)
   - [Standard Layout](#standard-layout)
   - [Separation of Concerns](#separation-of-concerns)
   - [Dependency Injection](#dependency-injection)

6. [Production Readiness](#6-production-readiness)
   - [Configuration Management](#configuration-management)
   - [Graceful Shutdown](#graceful-shutdown)
   - [Error Handling and Logging](#error-handling-and-logging)

---

## 1. Go Language Fundamentals

### Types and Data Structures

Go is a statically typed language with a focus on simplicity. Types are explicitly declared and include basic types (int, string, bool), composite types (arrays, slices, maps, structs), and reference types (pointers, channels, functions, interfaces).

#### Structs

Structs are collections of fields, similar to classes in other languages but without inheritance.

```go
// From our project: models/consultant.go
type Consultant struct {
    ID       int    `json:"id"`
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    SkillIDs []int  `json:"skill_ids"`
}
```

**Key Points:**
- Structs define data structures with named fields
- Field tags (like `json:"id"`) provide metadata for reflection-based operations
- Exported fields start with uppercase letters (public), unexported with lowercase (private)
- No inheritance; Go uses composition over inheritance
- No constructors; use factory functions instead

#### Maps

Maps are key-value stores, similar to hash tables or dictionaries in other languages.

```go
// From our in-memory store implementation
type Store struct {
    consultants map[int]models.Consultant
    skills      map[int]models.Skill
    mutex       sync.RWMutex
}
```

**Key Points:**
- Maps are reference types (passing a map passes a reference)
- Maps are not thread-safe; concurrent access requires synchronization
- Maps must be initialized with `make()` before use
- The zero value of a map is `nil`

#### Slices

Slices are dynamic arrays that offer a flexible view into an underlying array.

```go
// Get all consultants
func (db *PostgresDB) GetAllConsultants() ([]models.Consultant, error) {
    // Collect consultants
    var consultants []models.Consultant
    for rows.Next() {
        var c models.Consultant
        if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
            return nil, err
        }
        consultants = append(consultants, c)
    }
    
    return consultants, nil
}
```

**Key Points:**
- Slices are reference types with three components: pointer, length, capacity
- `append()` grows a slice dynamically when needed
- Slices can be sliced with `s[2:5]` notation
- Slices are commonly used instead of arrays due to their flexibility

### Packages and Imports

Go organizes code into packages, which are directories containing Go source files.

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/yourusername/consultancy-api/database"
    "github.com/yourusername/consultancy-api/models"
)
```

**Key Points:**
- Every Go file begins with a `package` declaration
- The `main` package is special - it defines an executable program
- Import paths use the full module path
- External packages are downloaded and versioned using Go modules
- Package names are lowercase, single-word identifiers
- Exported identifiers start with an uppercase letter

### Error Handling Philosophy

Go handles errors as values, not exceptions. Functions return errors as values that must be explicitly checked.

```go
// Get a consultant by ID
func (db *PostgresDB) GetConsultant(id int) (models.Consultant, error) {
    // Get consultant
    var consultant models.Consultant
    err := db.db.QueryRowContext(
        ctx,
        "SELECT id, name, email FROM consultants WHERE id = $1",
        id,
    ).Scan(&consultant.ID, &consultant.Name, &consultant.Email)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.Consultant{}, fmt.Errorf("consultant with id %d not found", id)
        }
        return models.Consultant{}, err
    }

    return consultant, nil
}
```

**Key Points:**
- Functions can return multiple values, commonly including an error
- Errors are checked immediately after the operation
- The `error` interface requires only an `Error() string` method
- Use `errors.New()` or `fmt.Errorf()` to create simple errors
- `errors.Is()` and `errors.As()` help with error checking and type assertions
- Zero values are commonly returned with errors

### Functions and Methods

Functions are first-class citizens in Go. Methods are functions associated with a specific type.

```go
// Function
func NewConsultantHandler(db *database.PostgresDB) *ConsultantHandler {
    return &ConsultantHandler{
        db: db,
    }
}

// Method with pointer receiver
func (h *ConsultantHandler) Create(w http.ResponseWriter, r *http.Request) {
    var consultant models.Consultant
    
    if err := json.NewDecoder(r.Body).Decode(&consultant); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    
    // Rest of method...
}
```

**Key Points:**
- Functions have zero or more parameters and zero or more return values
- Parameters are passed by value (copied)
- Functions can return multiple values
- Methods are functions with a receiver argument
- Pointer receivers can modify the receiver, value receivers cannot
- Use pointer receivers for methods that modify the receiver or for large structs
- Use value receivers for immutable operations or small structs

## 2. Concurrency in Go

### Goroutines

Goroutines are lightweight threads managed by the Go runtime. They allow concurrent execution with minimal resources.

```go
// Start HTTP server in a goroutine
go func() {
    log.Printf("Starting server on %s", srv.Addr)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        serverErrors <- err
    }
}()
```

**Key Points:**
- Created by adding the `go` keyword before a function call
- Much lighter than OS threads (starting with ~2KB of stack)
- The Go runtime multiplexes goroutines onto OS threads
- No direct access or reference to a goroutine
- Communicate through channels, not shared memory

### Channels

Channels are typed conduits for communication between goroutines, providing synchronized, thread-safe data exchange.

```go
// Channel for server errors
serverErrors := make(chan error, 1)

// Wait for interrupt signal or server error
select {
case err := <-serverErrors:
    log.Fatalf("Server error: %v", err)
case <-stop:
    log.Println("Shutting down server...")
    // Graceful shutdown...
}
```

**Key Points:**
- Created with `make(chan Type, [capacity])`
- Unbuffered channels (capacity 0) block until both sender and receiver are ready
- Buffered channels block only when the buffer is full
- Send with `ch <- value`, receive with `value := <-ch`
- Close with `close(ch)` when no more values will be sent
- Receivers can check for closed channels with `value, ok := <-ch`
- The `range` keyword iterates over values from a channel until it's closed

### WaitGroups

WaitGroups are a synchronization mechanism for waiting for a collection of goroutines to finish.

```go
// Example from a concurrent data fetching function
func fetchRelatedData(ids []int) []Data {
    var wg sync.WaitGroup
    results := make(chan Data, len(ids))
    
    for _, id := range ids {
        wg.Add(1)  // Increment counter
        go func(id int) {
            defer wg.Done()  // Decrement counter when done
            // Fetch data and send to results channel
            data := fetchSingleItem(id)
            results <- data
        }(id)
    }
    
    // Wait for all goroutines to finish, then close channel
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var items []Data
    for result := range results {
        items = append(items, result)
    }
    
    return items
}
```

**Key Points:**
- `Add(n)` increments the WaitGroup counter by n
- `Done()` decrements the counter by 1
- `Wait()` blocks until the counter is 0
- Typically used with `defer wg.Done()` to ensure counter is decremented
- Often combined with channels to collect results from goroutines

### Mutex

Mutexes (mutual exclusion locks) protect shared resources from concurrent access, preventing data races.

```go
// From our in-memory store
type Store struct {
    consultants map[int]models.Consultant
    mutex       sync.RWMutex  // Reader/Writer mutex
}

// Read operation with read lock
func (s *Store) GetConsultant(id int) (models.Consultant, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    consultant, exists := s.consultants[id]
    if !exists {
        return models.Consultant{}, fmt.Errorf("consultant not found")
    }
    
    return consultant, nil
}

// Write operation with write lock
func (s *Store) CreateConsultant(consultant models.Consultant) models.Consultant {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    // Modify shared data
    consultant.ID = s.nextConsultantID
    s.nextConsultantID++
    s.consultants[consultant.ID] = consultant
    
    return consultant
}
```

**Key Points:**
- `Mutex` provides exclusive access with `Lock()` and `Unlock()`
- `RWMutex` distinguishes between read and write operations
- Multiple readers can acquire a read lock simultaneously
- Only one writer can hold a write lock, blocking all readers
- Always use `defer` to ensure unlocking even when panics occur
- Maps and slices are not thread-safe and require mutexes for concurrent access

### Context Package

The context package provides a way to carry deadlines, cancellation signals, and request-scoped values across API boundaries and between goroutines.

```go
// Database operation with context timeout
func (db *PostgresDB) GetConsultant(id int) (models.Consultant, error) {
    // Use a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Use context in database query
    var consultant models.Consultant
    err := db.db.QueryRowContext(
        ctx,
        "SELECT id, name, email FROM consultants WHERE id = $1",
        id,
    ).Scan(&consultant.ID, &consultant.Name, &consultant.Email)
    
    // Rest of function...
}

// Graceful server shutdown with context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("Server forced to shutdown: %v", err)
}
```

**Key Points:**
- `context.Background()` is the root context, typically starting point
- `context.WithTimeout()` creates a context that will be canceled after a duration
- `context.WithCancel()` creates a context with a cancel function
- `context.WithValue()` stores key-value pairs in a context
- Always pass context as the first parameter of a function
- Always cancel contexts when they're no longer needed
- Never store contexts in structs; pass them explicitly
- Context values should be used for request-scoped data like trace IDs, not for passing optional parameters

## 3. Building REST APIs

### HTTP Package

The standard library's `net/http` package provides HTTP client and server implementations.

```go
// Basic HTTP server
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
})
log.Fatal(http.ListenAndServe(":8080", nil))
```

**Key Points:**
- `http.ListenAndServe` starts an HTTP server
- `http.Handler` is an interface with a `ServeHTTP(ResponseWriter, *Request)` method
- `http.HandlerFunc` converts a function to a Handler
- `http.ResponseWriter` is an interface for writing HTTP responses
- `http.Request` represents an HTTP request
- The default `nil` handler uses the global `http.DefaultServeMux`

### Router/Mux

Routers or multiplexers direct HTTP requests to appropriate handlers based on the URL path and method.

```go
// Initialize router
r := mux.NewRouter()

// Apply middleware
r.Use(loggingMiddleware)

// API routes
apiRouter := r.PathPrefix("/api").Subrouter()

// Consultant routes with path variables and HTTP methods
apiRouter.HandleFunc("/consultants", consultantHandler.GetAll).Methods("GET")
apiRouter.HandleFunc("/consultants/{id:[0-9]+}", consultantHandler.Get).Methods("GET")
apiRouter.HandleFunc("/consultants", consultantHandler.Create).Methods("POST")
apiRouter.HandleFunc("/consultants/{id:[0-9]+}", consultantHandler.Update).Methods("PUT")
apiRouter.HandleFunc("/consultants/{id:[0-9]+}", consultantHandler.Delete).Methods("DELETE")
```

**Key Points:**
- The standard library provides a basic router with `http.ServeMux`
- Third-party routers like Gorilla Mux offer more features
- Routes can include path variables (`{id}`)
- Routes can be restricted to specific HTTP methods
- Subrouters can group routes with common prefixes or middleware
- Path variables are extracted with `mux.Vars(r)`
- Regular expressions can constrain path variables (`{id:[0-9]+}`)

### Handlers

Handlers process HTTP requests and generate responses.

```go
// Handler type
type ConsultantHandler struct {
    db *database.PostgresDB
}

// Handler method
func (h *ConsultantHandler) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid consultant ID", http.StatusBadRequest)
        return
    }
    
    consultant, err := h.db.GetConsultant(id)
    if err != nil {
        // Check if it's a not found error
        if err.Error() == "consultant with id "+strconv.Itoa(id)+" not found" {
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            http.Error(w, "Failed to get consultant: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(consultant)
}
```

**Key Points:**
- Handlers implement the `http.Handler` interface or use `http.HandlerFunc`
- `http.ResponseWriter` writes the response headers, status code, and body
- `http.Request` provides access to request details (method, URL, headers, body)
- Error responses use `http.Error` with appropriate status codes
- Handler methods often access dependencies through closures or struct fields
- Handlers should be stateless and safe for concurrent use

### Middleware

Middleware intercepts HTTP requests/responses to add cross-cutting concerns like logging, authentication, or request validation.

```go
// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Log the request
        log.Printf(
            "%s %s %s",
            r.Method,
            r.RequestURI,
            time.Since(start),
        )
    })
}

// Apply middleware
r.Use(loggingMiddleware)
```

**Key Points:**
- Middleware takes a handler and returns a new handler
- Middleware can execute code before and after the next handler
- Middleware is typically applied in an "onion" pattern (first middleware applied is the last to finish)
- Middleware can be applied globally to all routes or to specific routes
- Middleware can short-circuit the request handling if needed (e.g., for authentication)
- Common middleware functions include logging, authentication, CORS, panic recovery

### Request/Response Processing

Go provides built-in functions for processing HTTP requests and responses, especially for JSON handling.

```go
// Decode JSON request body
var consultant models.Consultant
if err := json.NewDecoder(r.Body).Decode(&consultant); err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
}

// Validate required fields
if consultant.Name == "" || consultant.Email == "" {
    http.Error(w, "Name and email are required", http.StatusBadRequest)
    return
}

// Process the request...

// Send JSON response
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(createdConsultant)
```

**Key Points:**
- `json.NewDecoder(r.Body).Decode(&v)` decodes JSON to a Go struct
- `json.NewEncoder(w).Encode(v)` encodes a Go struct to JSON
- `r.URL.Query()` retrieves URL query parameters
- `mux.Vars(r)` retrieves path variables (with Gorilla Mux)
- `r.FormValue()` retrieves form values
- `w.Header().Set()` sets response headers
- `w.WriteHeader()` sets the response status code
- `http.Error()` sends an error response with a status code

## 4. Database Integration

### SQL Package

Go's `database/sql` package provides a generic interface around SQL databases, with specific drivers for different databases.

```go
import (
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver
)

// Open connection
db, err := sql.Open("postgres", connStr)
if err != nil {
    return nil, err
}
```

**Key Points:**
- `database/sql` is the standard interface for SQL databases
- Specific drivers (like `github.com/lib/pq` for PostgreSQL) implement the interface
- The driver is imported with the blank identifier `_` for its side effects
- `sql.Open` returns a handle to the database, not a connection
- The database handle manages a pool of connections
- `db.Ping()` verifies a connection can be established

### Connection Pooling

The `sql.DB` type manages a pool of database connections for optimal performance.

```go
// Configure connection pool
db.SetMaxOpenConns(25)      // Maximum number of open connections
db.SetMaxIdleConns(5)       // Maximum number of idle connections
db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime
```

**Key Points:**
- Connection pooling reuses connections to reduce overhead
- `SetMaxOpenConns` limits the total number of connections
- `SetMaxIdleConns` limits idle connections kept in the pool
- `SetConnMaxLifetime` sets the maximum time a connection can be reused
- Connections are automatically created and returned to the pool
- The pool handles connection failures and reconnects

### Transactions

Transactions group multiple operations into an atomic unit that either all succeed or all fail.

```go
// Begin a transaction
tx, err := db.db.BeginTx(ctx, nil)
if err != nil {
    return models.Consultant{}, err
}
defer tx.Rollback() // Will be ignored if transaction is committed

// Perform multiple operations
// 1. Insert consultant
err = tx.QueryRowContext(
    ctx,
    "INSERT INTO consultants (name, email) VALUES ($1, $2) RETURNING id",
    consultant.Name, consultant.Email,
).Scan(&consultant.ID)

if err != nil {
    return models.Consultant{}, err
}

// 2. Insert consultant skills
for _, skillID := range consultant.SkillIDs {
    _, err := tx.ExecContext(
        ctx,
        "INSERT INTO consultant_skills (consultant_id, skill_id) VALUES ($1, $2)",
        consultant.ID, skillID,
    )
    if err != nil {
        return models.Consultant{}, err
    }
}

// Commit transaction
if err := tx.Commit(); err != nil {
    return models.Consultant{}, err
}
```

**Key Points:**
- Transactions ensure atomicity, consistency, isolation, and durability (ACID)
- `BeginTx` starts a transaction with a context for timeout/cancellation
- `Commit` finalizes all changes
- `Rollback` abandons all changes
- `defer tx.Rollback()` is a safety net for aborted transactions
- Use transactions for operations that must succeed or fail as a unit
- Avoid long-running transactions that could block other operations

### Repository Pattern

The repository pattern abstracts data access logic from business logic, making code more maintainable and testable.

```go
// Database type (repository)
type PostgresDB struct {
    db *sql.DB
}

// Repository methods
func (db *PostgresDB) GetConsultant(id int) (models.Consultant, error) {
    // Implementation...
}

func (db *PostgresDB) CreateConsultant(consultant models.Consultant) (models.Consultant, error) {
    // Implementation...
}

// Handler uses repository
type ConsultantHandler struct {
    db *database.PostgresDB  // Repository dependency
}

func (h *ConsultantHandler) Get(w http.ResponseWriter, r *http.Request) {
    // Use repository
    consultant, err := h.db.GetConsultant(id)
    // Handle response...
}
```

**Key Points:**
- Repositories encapsulate data access logic
- Repository methods map to business operations, not raw CRUD
- Handlers depend on repositories, not on database details
- This separation enables unit testing with mock repositories
- Each model typically has its own repository
- Repositories can implement transaction management
- Repository interfaces can have multiple implementations (PostgreSQL, in-memory, etc.)

## 5. Project Organization

### Standard Layout

Go projects typically follow a standard layout that organizes code in a consistent way.

```
consultancy-api/
├── main.go              # Entry point
├── handlers/            # HTTP request handlers
│   ├── consultants.go
│   └── skills.go
├── models/              # Data structures
│   ├── consultant.go
│   └── skill.go
└── database/            # Database layer
    └── postgres.go
```

**Key Points:**
- `main.go` contains the application entry point and setup
- Each package resides in its own directory
- Package names match their directory names
- Related functionality is grouped into packages
- Keep packages small and focused on a single responsibility
- Avoid circular dependencies between packages
- Keep internal implementation details unexported

### Separation of Concerns

Go projects typically separate code into distinct layers, each with a specific responsibility.

```
Models      - Data structures and domain entities
Database    - Data access and storage
Handlers    - HTTP request handling and routing
Main        - Application setup and configuration
```

**Key Points:**
- Models define the shape of data and domain entities
- Database layer handles data persistence and retrieval
- Handlers manage HTTP requests and responses
- Main coordinates the application components
- Each layer depends only on the layers below it
- Lower layers should not depend on higher layers
- This separation enables easier testing and maintenance

### Dependency Injection

Go typically uses explicit dependency injection to manage component dependencies.

```go
// Database dependency
db, err := database.New(dbConfig)
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}

// Inject database into handlers
consultantHandler := handlers.NewConsultantHandler(db)
skillHandler := handlers.NewSkillHandler(db)

// Use handlers
r.HandleFunc("/consultants", consultantHandler.GetAll).Methods("GET")
```

**Key Points:**
- Dependencies are passed explicitly through constructors
- Components receive their dependencies, not create them
- This approach enables easier testing with mock dependencies
- No dependency injection frameworks needed; use simple constructors
- Dependencies are typically stored in struct fields
- Factory functions create and configure components

## 6. Production Readiness

### Configuration Management

Production applications need flexible configuration for different environments.

```go
// Load environment variables from .env file
if err := godotenv.Load(); err != nil {
    log.Println("No .env file found, using environment variables")
}

// Database configuration
dbConfig := database.Config{
    Host:     getEnv("DB_HOST", "localhost"),
    Port:     getEnvAsInt("DB_PORT", 5432),
    User:     getEnv("DB_USER", "postgres"),
    Password: getEnv("DB_PASSWORD", "postgres"),
    DBName:   getEnv("DB_NAME", "consultancy"),
    SSLMode:  getEnv("DB_SSLMODE", "disable"),
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
```

**Key Points:**
- Environment variables are a common configuration source
- `.env` files can provide local development configuration
- Default values ensure the application can run without explicit configuration
- Configuration should be validated at startup
- Sensitive values (passwords, API keys) should be handled securely
- Configuration should be environment-specific (dev, staging, production)

### Graceful Shutdown

Production services should shut down gracefully to avoid disrupting clients.

```go
// Channel for OS signals
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt)

// Wait for interrupt signal
<-stop
log.Println("Shutting down server...")

// Create a deadline for server shutdown
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Attempt graceful shutdown
if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("Server forced to shutdown: %v", err)
}

log.Println("Server gracefully stopped")
```

**Key Points:**
- Graceful shutdown allows in-flight requests to complete
- OS signals (like SIGINT from Ctrl+C) trigger shutdown
- `context.WithTimeout` sets a maximum shutdown time
- `srv.Shutdown` stops accepting new connections but waits for existing requests
- Resources should be closed in the correct order (server first, then database)
- Log messages should indicate shutdown progress

### Error Handling and Logging

Robust error handling and logging are essential for production applications.

```go
// Handler error handling
func (h *ConsultantHandler) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Printf("Invalid consultant ID: %v", err)
        http.Error(w, "Invalid consultant ID", http.StatusBadRequest)
        return
    }
    
    consultant, err := h.db.GetConsultant(id)
    if err != nil {
        // Check if it's a not found error
        if err.Error() == "consultant with id "+strconv.Itoa(id)+" not found" {
            log.Printf("Consultant not found: %d", id)
            http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            log.Printf("Database error: %v", err)
            http.Error(w, "Failed to get consultant: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }
    
    // Success path...
}
```

**Key Points:**
- Log all errors with context (request details, operation type)
- Map errors to appropriate HTTP status codes
- Don't expose internal error details to clients in production
- Use structured logging for easier parsing and analysis
- Log important events (startup, shutdown, configuration)
- Consider using a logging library for advanced features
- Include trace IDs for request correlation
- Use different log levels (debug, info, warn, error) appropriately

---# go-service-api
