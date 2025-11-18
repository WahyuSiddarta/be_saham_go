# Middleware Package Documentation

This package contains all middleware logic for the Go API application, providing a centralized location for request processing, security, and monitoring.

## Structure

````
middleware/
├── middleware.go     # Main package file with shared logger
├── setup.go         # Middleware configuration and setup
├── cors.go          # CORS (Cross-Origin Resource Sharing) middleware
├── ratelimiter.go   # Rate limiting middleware
├── auth.go          # Authentication middleware (placeholder)
├── logging.go       # Request logging middleware
└── recover.go       # Panic recovery middleware
```## Available Middleware

### 1. **CORS Middleware** (`cors.go`)

Handles Cross-Origin Resource Sharing for frontend integration.

**Functions:**

- `ConfigureCORS()` - Returns configured CORS middleware based on environment settings
- `LogCORSStatus()` - Logs current CORS configuration

**Configuration:**

```bash
CORS_ENABLED=true
CORS_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:4173
````

### 2. **Rate Limiting Middleware** (`ratelimiter.go`)

Protects against abuse and ensures fair usage of API resources.

**Functions:**

- `NewRateLimiter()` - Creates a new rate limiter instance
- `Middleware()` - Returns rate limiting middleware function
- `LogRateLimitStatus()` - Logs current rate limiting configuration
- `CleanupVisitors()` - Removes old visitor entries to prevent memory leaks

**Configuration:**

```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200
```

**Features:**

- Per-IP rate limiting using token bucket algorithm
- Configurable requests per second and burst size
- Automatic cleanup to prevent memory leaks
- Structured logging of violations

### 3. **Authentication Middleware** (`auth.go`)

Provides authentication and authorization (placeholder for future implementation).

**Functions:**

- `AuthMiddleware()` - Basic auth middleware (placeholder)
- `RequireAuth()` - Returns unauthorized response (placeholder)

**Status:** Currently returns structured error responses indicating auth is not implemented.

### 4. **Request Logging Middleware** (`logging.go`)

Logs all HTTP requests with detailed information.

**Functions:**

- `RequestLogger()` - Logs all requests with timing and status information
- `HealthCheckLogger()` - Lighter logging specifically for health check endpoints

**Logged Information:**

- HTTP method and path
- Remote IP address
- Response status code
- Request duration
- User agent
- Error details (if any)

### 5. **Panic Recovery Middleware** (`recover.go`)

Gracefully handles panics that occur during request processing.

**Functions:**

- `Recover()` - Returns recovery middleware with default configuration
- `RecoverWithConfig()` - Returns recovery middleware with custom configuration

**Features:**

- Captures and logs panic stack traces
- Returns proper HTTP 500 error responses
- Prevents server crashes from unhandled panics
- Configurable stack trace size and options
- Telegram notification support (commented out, can be enabled)

**Configuration:**

```go
type RecoverConfig struct {
    StackTraceSize                 int  // Memory allocated for stack trace
    PrintStackTraceOfAllGoroutines bool // Include all goroutines in trace
    ErrorHandler                   func(c echo.Context, err error) // Custom error handler
}
```

### 6. **Setup Middleware** (`setup.go`)

Centralizes middleware configuration and application.

**Functions:**

- `SetupGlobalMiddleware(e *echo.Echo)` - Configures global middleware for entire application
- `SetupAPIMiddleware(apiGroup *echo.Group)` - Configures middleware specific to API routes

## Usage

### In main.go:

```go
import "github.com/WahyuSiddarta/be_saham_go/middleware"

// Setup global middleware
middleware.SetupGlobalMiddleware(echoInstance)
```

### In router setup:

```go
// Setup API-specific middleware
apiGroup := router.Group("/api")
middleware.SetupAPIMiddleware(apiGroup)
```

### Individual middleware usage:

```go
// Apply specific middleware to route groups
protectedGroup.Use(middleware.RequireAuth())
```

## Middleware Order

The middleware is applied in this order:

1. **Global Level (all routes):**

   - Panic Recovery
   - Request Logging
   - CORS

2. **API Level (/api/\* routes):**

   - Rate Limiting

3. **Route Group Level:**
   - Authentication (for protected routes)

## Configuration

All middleware can be enabled/disabled via environment variables:

```bash
# CORS
CORS_ENABLED=true
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200
```

## Logging

All middleware integrates with the application's structured logging system:

- **Info logs:** Configuration status, successful operations
- **Warning logs:** Rate limit violations, auth failures
- **Error logs:** Request processing errors
- **Debug logs:** Health check requests (to reduce log noise)

## Extension Points

### Adding New Middleware

1. Create a new file in the middleware package (e.g., `security.go`)
2. Implement the middleware function returning `echo.MiddlewareFunc`
3. Add configuration to `config/config.go` if needed
4. Add to setup functions in `setup.go`
5. Update logger distribution if needed

### Example New Middleware:

```go
// security.go
package middleware

func SecurityHeaders() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Response().Header().Set("X-Frame-Options", "DENY")
            c.Response().Header().Set("X-Content-Type-Options", "nosniff")
            return next(c)
        }
    }
}
```

## Best Practices

1. **Order Matters:** Apply middleware in logical order (logging first, auth last)
2. **Configuration:** Make middleware configurable via environment variables
3. **Logging:** Use structured logging with appropriate log levels
4. **Performance:** Consider performance impact of middleware order
5. **Error Handling:** Return consistent error responses
6. **Memory Management:** Implement cleanup for stateful middleware (like rate limiting)

## Testing Middleware

Test middleware functionality:

```bash
# Test rate limiting
for i in {1..150}; do curl http://localhost:3000/api/public/health; done

# Test CORS
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: POST" \
     -X OPTIONS \
     http://localhost:3000/api/public/health

# Test authentication (should return auth not implemented)
curl http://localhost:3000/api/protected/profile
```
