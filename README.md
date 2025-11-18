# be_saham_go

Backend API service for Indonesian stock market data built with Go and Echo framework.

## Features

- RESTful API for stock market data
- Real-time stock price information
- Stock price history
- CORS enabled for frontend integration
- Structured logging
- Health check endpoints

## Tech Stack

- **Framework**: Echo v4
- **Language**: Go 1.21+
- **Architecture**: Clean architecture with handlers, models, and routes

## API Endpoints

### Health Check

- `GET /` - Welcome message
- `GET /health` - Health status

### Stock API

- `GET /api/v1/stocks` - Get all stocks
- `GET /api/v1/stocks/:symbol` - Get specific stock by symbol
- `GET /api/v1/stocks/:symbol/prices` - Get price history for a stock

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/WahyuSiddarta/be_saham_go.git
cd be_saham_go
```

2. Install dependencies:

```bash
go mod tidy
```

3. Copy environment variables:

```bash
cp .env.example .env
```

4. Run the application:

```bash
# Direct Go command
go run main.go

# Or use the convenience script
./run.sh
```

The server will start on `http://localhost:3000` (or port specified in .env)

### Development

```bash
# Build the application
go build -o bin/be_saham_go .

# Run tests
go test ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy

# Install air for hot reload (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reload (if air is installed)
air
```

## Project Structure

```
be_saham_go/
├── api/              # API handlers and logic
├── config/           # Configuration management
├── db/               # Database connections and utilities
├── helper/           # Helper functions and utilities
├── middleware/       # HTTP middleware (CORS, Rate limiting, Auth, Logging)
├── models/           # Data models and structures
├── router/           # Route definitions and setup
├── logger/           # Logging configuration and setup
├── log/              # Log files directory
├── main.go           # Application entry point
├── go.mod            # Go module file
└── README.md         # Project documentation
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
