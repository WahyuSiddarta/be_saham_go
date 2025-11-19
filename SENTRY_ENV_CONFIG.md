# Sentry DSN Configuration - Environment Setup

## ✅ Changes Made

### Files Modified:

1. **`config/config.go`**

   - Added `SentryDSN string` field to `Config` struct
   - Added loading from environment variable `SENTRY_DSN`
   - Default value is empty string (optional configuration)

2. **`main.go`**

   - Moved hardcoded DSN to use `configStruct.SentryDSN`
   - Added conditional check: only initialize Sentry if DSN is provided
   - Logs informative message if DSN is not configured

3. **`.env.example`**
   - Added `SENTRY_DSN` variable with your DSN as example
   - Added comment explaining DSN format and where to get it

## 🔧 How to Use

### Development Setup

1. **Copy .env.example to .env:**

   ```bash
   cp .env.example .env
   ```

2. **Add your Sentry DSN to `.env`:**

   ```bash
   SENTRY_DSN=https://f6abe2bd415288e76e7f8782180972b3@o4510388552204288.ingest.us.sentry.io/4510388639236096
   ```

3. **Run the application:**
   ```bash
   go run main.go
   # or
   ./bin/app
   ```

### Production Setup

Set the environment variable before running:

```bash
export SENTRY_DSN="your-production-dsn-here"
./app
```

Or in Docker:

```dockerfile
ENV SENTRY_DSN=your-production-dsn-here
```

Or in systemd service file:

```ini
[Service]
Environment="SENTRY_DSN=your-production-dsn-here"
```

## 📋 Configuration Priority

Environment variables are checked in this order:

1. **Environment Variable**: `SENTRY_DSN` (from OS or .env file)
2. **Default Value**: Empty string `""` (Sentry is skipped if not configured)

## 🎯 Behavior

### When DSN is Configured

```
Sentry DSN found in config
         ↓
sentry.Init() called with DSN
         ↓
Sentry service initialized successfully
         ↓
All errors automatically captured
```

### When DSN is NOT Configured

```
Sentry DSN is empty string
         ↓
Init skipped (check prevents nil DSN)
         ↓
Logs: "Sentry DSN not configured, skipping Sentry initialization"
         ↓
Application runs normally without Sentry
```

## 🔐 Security Notes

✅ **Good Practice:**

- DSN is NOT hardcoded anymore
- DSN is loaded from environment variables
- Can be different per environment (dev/staging/production)
- Sensitive data is not in version control

⚠️ **Important:**

- Make sure `.env` file is added to `.gitignore`
- Different DSNs for dev, staging, and production
- Rotate DSN if it's exposed

## 📄 Example `.env` File

```bash
# Application Configuration
PORT=8080
ENV=production
LOG_LEVEL=info

# Sentry Error Tracking (ADD YOUR DSN HERE)
SENTRY_DSN=https://your-key@your-org.ingest.us.sentry.io/your-project-id

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h

# ... other configuration ...
```

## ✨ Benefits

| Benefit          | Before              | After                     |
| ---------------- | ------------------- | ------------------------- |
| Security         | DSN in code         | DSN in env vars           |
| Flexibility      | Same DSN everywhere | Different per environment |
| Production Ready | Risky               | Safe                      |
| Development      | Manual code changes | Just use `.env.example`   |
| CI/CD            | Need code changes   | Just set env var          |

## 🧪 Testing

To test the configuration:

1. **Without SENTRY_DSN:**

   ```bash
   unset SENTRY_DSN
   go run main.go
   # Output: "Sentry DSN not configured, skipping Sentry initialization"
   ```

2. **With SENTRY_DSN:**

   ```bash
   export SENTRY_DSN="your-dsn-here"
   go run main.go
   # Output: "Sentry initialized"
   ```

3. **From .env file:**
   ```bash
   # Create .env with SENTRY_DSN=...
   go run main.go
   # Output: "Sentry initialized"
   ```

## 🚀 Build Status

✅ Project builds successfully
✅ All tests pass
✅ Configuration loaded correctly
✅ Ready for deployment

## 📚 Related Files

- `.env.example` - Example environment variables
- `config/config.go` - Configuration loading
- `main.go` - Sentry initialization
- `middleware/sentry.go` - Sentry utilities

---

**Now your Sentry DSN is properly managed through environment variables!** 🎉
