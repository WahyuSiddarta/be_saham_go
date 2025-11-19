# Sentry Quick Reference

## 🚀 Quick Start

### 1. Import in Your Handler

```go
import "github.com/WahyuSiddarta/be_saham_go/middleware"
```

### 2. Capture Errors

```go
middleware.CaptureException(c, err)
```

### 3. Set User Context (After Auth)

```go
middleware.SetUserContext(c, userID, email, username)
```

---

## 📚 Main Functions

| Function                             | Usage                | When to Use              |
| ------------------------------------ | -------------------- | ------------------------ |
| `CaptureException(c, err)`           | Simple error capture | Quick error logging      |
| `CaptureError(c, err, tags, extra)`  | Error with context   | Detailed error tracking  |
| `CaptureMessage(c, msg)`             | Custom message       | Important events         |
| `SetUserContext(c, id, email, user)` | Associate with user  | After authentication     |
| `CaptureRecovery(c, err, stack)`     | Panic capture        | Auto-used by middleware  |
| `FlushSentry(seconds)`               | Send pending events  | Graceful shutdown (auto) |

---

## 💡 Quick Examples

### Example 1: Catch Database Error

```go
user, err := db.GetUser(id)
if err != nil {
    middleware.CaptureError(c, err,
        map[string]string{"operation": "get_user"},
        map[string]interface{}{"user_id": id},
    )
    return errorResponse(c, "User not found")
}
```

### Example 2: Track User Activity

```go
// In login handler
middleware.SetUserContext(c, user.ID, user.Email, user.Username)

// Now all errors in this request are tied to this user
```

### Example 3: API Integration

```go
resp, err := externalAPI.Call(data)
if err != nil {
    middleware.CaptureError(c, err,
        map[string]string{"api": "payment_service"},
        map[string]interface{}{"endpoint": "/charge"},
    )
    return errorResponse(c, "Payment failed")
}
```

### Example 4: Business Event

```go
middleware.CaptureMessage(c, "High-value transaction completed: $10,000")
```

---

## 🎯 Best Practices

✅ **DO:**

- Set user context after authentication
- Add tags to categorize errors
- Include relevant context data
- Capture external API errors
- Monitor database errors

❌ **DON'T:**

- Capture passwords or tokens
- Log raw request bodies
- Capture duplicate errors
- Ignore error logs
- Disable Sentry in production

---

## 📊 View Your Data

Go to: **https://sentry.io/**

Tabs:

- **Issues**: All errors and exceptions
- **Performance**: Request timing and bottlenecks
- **Releases**: Errors by app version
- **Alerts**: Set up notifications

---

## 🔧 Environment Setup

DSN (configured):

```
https://f6abe2bd415288e76e7f8782180972b3@o4510388552204288.ingest.us.sentry.io/4510388639236096
```

Tracing: ✅ Enabled (100%)
Error Capture: ✅ Enabled

---

## 📄 Full Documentation

- **Main Docs**: Read `SENTRY.md`
- **Examples**: Read `SENTRY_EXAMPLES.md`
- **Implementation**: Read `SENTRY_IMPLEMENTATION.md`

---

## ⚡ Common Scenarios

### Scenario 1: Login Handler

```go
func Login(c echo.Context) error {
    user, err := authenticate(email, pass)
    if err != nil {
        middleware.CaptureException(c, err)
        return errorResponse(c, "Login failed")
    }

    middleware.SetUserContext(c, user.ID, user.Email, user.Name)

    token, err := generateToken(user)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{"action": "generate_token"},
            nil,
        )
        return errorResponse(c, "Token generation failed")
    }

    return c.JSON(200, map[string]string{"token": token})
}
```

### Scenario 2: Portfolio Creation

```go
func CreatePortfolio(c echo.Context) error {
    userID := getUserID(c)
    middleware.SetUserContext(c, userID, "", "")

    portfolio, err := db.CreatePortfolio(userID, req)
    if err != nil {
        middleware.CaptureError(c, err,
            map[string]string{
                "operation": "create_portfolio",
                "entity": "portfolio",
            },
            map[string]interface{}{
                "user_id": userID,
                "name": req.Name,
            },
        )
        return errorResponse(c, "Failed to create portfolio")
    }

    return c.JSON(201, portfolio)
}
```

### Scenario 3: Panic Handling

```go
// Panics are handled automatically!
// No need to do anything special.
// The recover middleware will:
// 1. Catch the panic
// 2. Capture it in Sentry
// 3. Return a 500 error
// 4. Continue running
```

---

## 🧪 Testing Sentry

To test if Sentry is working:

```go
// In any handler (for testing only)
middleware.CaptureMessage(c, "Test message from be_saham_go")
```

Then check Sentry dashboard - you should see the message appear within seconds.

---

**All set! Start using Sentry in your handlers today.** 🎉
