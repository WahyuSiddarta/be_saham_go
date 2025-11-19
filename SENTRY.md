# Sentry Integration Documentation

This document describes the Sentry error tracking and performance monitoring implementation in the be_saham_go project.

## Overview

Sentry is integrated to:

- **Capture exceptions and errors** automatically
- **Track performance** with distributed tracing
- **Monitor panics** and critical failures
- **Manage user context** for better error attribution
- **Alert on issues** in real-time

## Configuration

### DSN (Data Source Name)

```
https://f6abe2bd415288e76e7f8782180972b3@o4510388552204288.ingest.us.sentry.io/4510388639236096
```

The DSN is configured in `main.go` during application initialization with:

- **EnableTracing**: true - Enable performance monitoring
- **TracesSampleRate**: 1.0 - Capture 100% of transactions

## Implementation Details

### 1. Middleware Integration

Sentry middleware is automatically integrated in the Echo request pipeline via `middleware/setup.go`:

```go
e.Use(sentryecho.New(sentryecho.Options{
    Repanic: true,  // Re-raise panics after capturing
}))
```

**Key Features:**

- Automatically captures HTTP errors (4xx, 5xx)
- Tracks request/response timing
- Captures request headers and body (when appropriate)
- Sets up user context from request data

### 2. Panic Recovery

Panics are captured in `middleware/recover.go`:

- Stack trace is captured with full details
- Panic information is sent to Sentry with level: FATAL
- Request context is preserved
- Application continues to run (graceful recovery)

### 3. Helper Functions

The `middleware/sentry.go` file provides utility functions for manual error tracking:

#### CaptureException

```go
CaptureException(c echo.Context, err error)
```

Captures an exception with request context.

#### CaptureError

```go
CaptureError(c echo.Context, err error, tags map[string]string, extra map[string]interface{})
```

Captures an error with custom tags and extra context.

#### CaptureMessage

```go
CaptureMessage(c echo.Context, message string)
```

Captures an informational message.

#### SetUserContext

```go
SetUserContext(c echo.Context, userID interface{}, userEmail string, username string)
```

Associates errors with specific users for better tracking.

#### CaptureRecovery

```go
CaptureRecovery(c echo.Context, err error, stackTrace string)
```

Captures panic recovery events.

#### FlushSentry

```go
FlushSentry(timeoutSeconds int) bool
```

Flushes pending events during graceful shutdown (called in `main.go`).

## Usage Examples

### Manual Error Capture in API Handlers

```go
// In your API handler
err := someOperation()
if err != nil {
    middleware.CaptureError(c, err,
        map[string]string{"operation": "create_portfolio"},
        map[string]interface{}{"user_portfolio": portfolioID},
    )
    return helper.ErrorResponse(c, http.StatusBadRequest, "Operation failed", nil)
}
```

### Setting User Context

```go
// After user authentication
userID := getUserID(c)
middleware.SetUserContext(c, userID, user.Email, user.Username)
```

### Capturing Specific Messages

```go
// For important events
middleware.CaptureMessage(c, "Unusual transaction detected")
```

## Automatic Tracking

The following are automatically captured by Sentry:

1. **HTTP Errors**: All 4xx and 5xx responses
2. **Exceptions**: Unhandled errors and panics
3. **Transactions**: Request/response timing
4. **Context**: Request method, URL, query parameters
5. **Performance**: Database queries, external API calls (with instrumentation)

## Viewing Events in Sentry

1. Go to [Sentry Dashboard](https://sentry.io/)
2. Log in with your account
3. Navigate to the be_saham_go project
4. View:
   - **Issues**: Grouped errors and exceptions
   - **Performance**: Transaction timings and bottlenecks
   - **Releases**: Track errors by application version
   - **Alerts**: Configure notifications for critical errors

## Environment Variables

The Sentry DSN is hardcoded in `main.go`. To make it configurable, add to `.env`:

```bash

```

Then update `main.go`:

```go
sentryDSN := configStruct.SentryDSN // Add this to config.go
if err := sentry.Init(sentry.ClientOptions{
    Dsn:              sentryDSN,
    EnableTracing:    true,
    TracesSampleRate: 1.0,
}); err != nil {
    Logger.Warn().Err(err).Msg("Sentry initialization failed")
}
```

## Best Practices

1. **Set user context early** after authentication to track user-specific errors
2. **Add custom tags** to categorize errors by feature or operation
3. **Include extra context** for debugging (user IDs, resource IDs, etc.)
4. **Avoid capturing sensitive data** in error messages or contexts
5. **Monitor release deployments** to track error trends across versions
6. **Configure alerts** for critical issues that need immediate attention

## Performance Monitoring

Sentry tracks:

- HTTP request duration
- Database query timing
- External API call latency
- Error rates by endpoint

Configure additional instrumentation as needed for specific operations.

## Dependencies

```
github.com/getsentry/sentry-go v0.38.0
github.com/getsentry/sentry-go/echo v0.38.0
```

## Shutdown Handling

Sentry is properly flushed on application shutdown in `main.go`:

```go
middleware.FlushSentry(5) // Wait max 5 seconds to flush pending events
```

This ensures all in-flight errors are captured before the application terminates.
