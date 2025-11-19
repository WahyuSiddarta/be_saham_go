/*
Package main - SENTRY IMPLEMENTATION EXAMPLES

This file contains documentation and examples of how to use Sentry integration
in your API handlers. Copy these patterns into your actual handler files.

=================================================================================
                          INTEGRATION CHECKLIST
=================================================================================

✓ 1. Dependencies added (sentry-go, sentry-go/echo)
✓ 2. Sentry initialized in main.go with DSN
✓ 3. Sentry middleware added to Echo setup in middleware/setup.go
✓ 4. Panic recovery middleware updated to capture panics in middleware/recover.go
✓ 5. Helper functions available in middleware/sentry.go
✓ 6. Graceful shutdown with Sentry flush in main.go

=================================================================================
                          USAGE IN YOUR API HANDLERS
=================================================================================

PATTERN 1: Simple Error Capture
------------------------------------
In any API handler file (e.g., api/api.auth.go):

    err := someOperation()
    if err != nil {
        middleware.CaptureException(c, err)
        return helper.ErrorResponse(c, http.StatusBadRequest, "Operation failed", nil)
    }


PATTERN 2: Error with Custom Tags and Context
------------------------------------
    result, err := complexOperation()
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{
                "operation": "complex_operation",
                "user_level": "premium",
            },
            map[string]interface{}{
                "user_id": userID,
                "resource_id": 123,
            },
        )
        return helper.ErrorResponse(c, http.StatusInternalServerError, "Operation failed", nil)
    }


PATTERN 3: Set User Context (do this after authentication)
------------------------------------
    userID := extractUserIDFromToken(token)
    middleware.SetUserContext(c, userID, user.Email, user.Username)
    // Now all errors in this request will be associated with this user


PATTERN 4: Capture Custom Messages
------------------------------------
    middleware.CaptureMessage(c, "Important business event occurred")


PATTERN 5: Database Error Handling
------------------------------------
    data, err := database.GetUser(userID)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{
                "error_type": "database",
                "operation": "select",
            },
            map[string]interface{}{
                "table": "users",
                "user_id": userID,
            },
        )
        return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user", nil)
    }


PATTERN 6: External API Error Handling
------------------------------------
    response, err := externalAPI.Charge(paymentData)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{
                "error_type": "external_api",
                "api_name": "payment_gateway",
            },
            map[string]interface{}{
                "endpoint": "https://api.payment.com/charge",
                "amount": paymentData.Amount,
            },
        )
        return helper.ErrorResponse(c, http.StatusPaymentRequired, "Payment processing failed", nil)
    }


=================================================================================
                          AUTOMATIC FEATURES
=================================================================================

The following are automatically captured by Sentry:

1. HTTP Errors
   - All 4xx and 5xx responses are automatically tracked
   - Request/response timing is recorded
   
2. Panics
   - Any panic in your handlers is caught by middleware/recover.go
   - Stack trace is captured with full details
   - Application continues running gracefully
   
3. Context Information
   - Request method, URL, query parameters
   - Request headers (sanitized)
   - Response status code
   - Execution time
   
4. User Information
   - User ID, email, username (when SetUserContext is called)
   - All subsequent errors are associated with that user


=================================================================================
                          SENTRY DASHBOARD
=================================================================================

Access your Sentry project at:
https://sentry.io/

View:
- Issues: Grouped errors and exceptions
- Performance: Transaction timings and bottlenecks
- Releases: Track errors by version
- Alerts: Configure notifications


=================================================================================
                          EXAMPLE IMPLEMENTATIONS
=================================================================================

EXAMPLE 1: Authentication Handler
------------------------------------
func (a *API) Login(c echo.Context) error {
    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        middleware.CaptureError(c, err,
            map[string]string{"action": "login", "step": "bind_json"},
            nil,
        )
        return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request", nil)
    }

    user, err := a.authenticateUser(req.Email, req.Password)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{"action": "login", "step": "authenticate"},
            map[string]interface{}{"email": req.Email},
        )
        return helper.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", nil)
    }

    // Set user context for this request
    middleware.SetUserContext(c, user.ID, user.Email, user.Username)

    token, err := generateJWT(user)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{"action": "login", "step": "jwt_generation"},
            map[string]interface{}{"user_id": user.ID},
        )
        return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", nil)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "token": token,
        "user":  user,
    })
}


EXAMPLE 2: Portfolio Creation Handler
------------------------------------
func (a *API) CreatePortfolio(c echo.Context) error {
    userID := extractUserID(c)
    middleware.SetUserContext(c, userID, "", "")

    var req CreatePortfolioRequest
    if err := c.BindJSON(&req); err != nil {
        middleware.CaptureException(c, err)
        return helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", nil)
    }

    portfolio, err := a.createUserPortfolio(userID, req)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{
                "operation": "create_portfolio",
                "user_level": getUserLevel(userID),
            },
            map[string]interface{}{
                "user_id": userID,
                "portfolio_name": req.Name,
            },
        )
        return helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create portfolio", nil)
    }

    return c.JSON(http.StatusCreated, portfolio)
}


EXAMPLE 3: Error Recovery (Automatic)
------------------------------------
// This is handled automatically by middleware/recover.go
// If any handler panics:

func (a *API) RiskyOperation(c echo.Context) error {
    userID := extractUserID(c)
    
    // If this panics, it will be:
    // 1. Caught by Recover() middleware
    // 2. Captured by CaptureRecovery() function
    // 3. Logged with full stack trace
    // 4. Returned as 500 Internal Server Error
    result := performRiskyOperation()
    
    return c.JSON(http.StatusOK, result)
}


=================================================================================
                          ENVIRONMENT CONFIGURATION
=================================================================================

Current Sentry DSN:
https://f6abe2bd415288e76e7f8782180972b3@o4510388552204288.ingest.us.sentry.io/4510388639236096

To make this configurable via environment variables:

1. Add to config/config.go:
   type Config struct {
       // ...existing fields...
       SentryDSN string
   }

2. Load from .env file:
   SENTRY_DSN=https://your-dsn-here

3. Update main.go initialization:
   if err := sentry.Init(sentry.ClientOptions{
       Dsn: configStruct.SentryDSN,
       EnableTracing: true,
       TracesSampleRate: 1.0,
   }); err != nil { ... }


=================================================================================
                          BEST PRACTICES
=================================================================================

1. ALWAYS Set User Context After Authentication
   middleware.SetUserContext(c, userID, email, username)

2. Add Custom Tags for Better Organization
   map[string]string{"feature": "payments", "action": "charge"}

3. Include Relevant Context Data
   map[string]interface{}{"amount": 100.00, "currency": "USD"}

4. Don't Capture Sensitive Data
   - Avoid passwords, API keys, tokens in error messages
   - Use placeholder values for sensitive IDs if needed

5. Monitor Error Trends
   - Review Sentry dashboard regularly
   - Set up alerts for critical errors
   - Track error rates across releases

6. Use Appropriate Error Levels
   - CaptureException(): For actual errors
   - CaptureMessage(): For informational messages
   - CaptureRecovery(): For panic recovery (automatic)

7. Disable in Development (Optional)
   if os.Getenv("ENV") == "production" {
       sentry.Init(...)
   }

=================================================================================
                          SHUTDOWN HANDLING
=================================================================================

Sentry is automatically flushed on graceful shutdown:
    middleware.FlushSentry(5) // Wait up to 5 seconds

This ensures all in-flight error events are sent before the application exits.

=================================================================================
*/

