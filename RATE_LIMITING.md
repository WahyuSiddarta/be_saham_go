# Rate Limiting Documentation

This application includes a built-in rate limiting system to protect against abuse and ensure fair usage of the API.

## Configuration

Rate limiting is configured through environment variables:

```bash
# Enable/disable rate limiting
RATE_LIMIT_ENABLED=true

# Requests per second allowed per IP
RATE_LIMIT_RPS=100

# Burst size (maximum number of requests that can be made in a burst)
RATE_LIMIT_BURST=200
```

## How it Works

1. **Per-IP Rate Limiting**: Each IP address gets its own rate limiter
2. **Token Bucket Algorithm**: Uses Go's `golang.org/x/time/rate` package
3. **Universal Application**: Applied to ALL routes under `/api/*`
4. **Memory Management**: Includes cleanup mechanisms to prevent memory leaks

## API Structure

All API endpoints are now consistently prefixed with `/api`:

```
/api/public/*     - Public endpoints (no authentication required)
/api/protected/*  - Protected endpoints (require authentication)
```

### Available Endpoints

#### Public Endpoints

- `GET /api/public/health` - Health check endpoint
- `GET /api/public/test` - Test endpoint

#### Protected Endpoints

- `GET /api/protected/profile` - User profile (placeholder)
- `GET /api/protected/stocks` - List stocks (placeholder)
- `GET /api/protected/stocks/:id` - Get specific stock (placeholder)

## Rate Limit Response

When rate limit is exceeded, the API returns:

```json
{
  "status": "ERR",
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Too many requests. Please try again later."
}
```

HTTP Status Code: `429 Too Many Requests`

## Testing Rate Limiting

You can test the rate limiting by making rapid requests to any `/api/*` endpoint:

```bash
# Start the server
go run main.go

# Make rapid requests to test rate limiting
for i in {1..150}; do curl http://localhost:8080/api/public/health; done
```

## Configuration Notes

- **Requests Per Second (RPS)**: Base rate of requests allowed per second
- **Burst Size**: Maximum number of requests that can be processed in a burst
- **Enabled Flag**: Allows you to completely disable rate limiting if needed
- **Development vs Production**: Adjust these values based on your environment

## Memory Considerations

The rate limiter keeps track of IP addresses in memory. In production environments with high traffic, consider:

1. Implementing Redis-based rate limiting for distributed systems
2. Regular cleanup of old visitor entries
3. Monitoring memory usage

## Security Benefits

- **DDoS Protection**: Prevents overwhelming the server with requests
- **Fair Usage**: Ensures all users get fair access to resources
- **Resource Protection**: Prevents abuse of expensive operations
- **Logging**: Logs rate limit violations for monitoring
