# 🚀 Ethereum Transaction Parser

[![Docker](https://img.shields.io/badge/Docker-Supported-blue)](https://www.docker.com/)

A service that monitors the Ethereum blockchain for transactions and provides notifications for subscribed addresses. This service allows users to subscribe to Ethereum addresses and query their transactions through a simple HTTP API.

## 📖 Table of Contents
- [Architecture](#architecture)
- [Key Design Decisions](#-key-design-decisions)
- [Prerequisites](#-prerequisites)
- [Getting Started](#-getting-started)
- [Running with Docker](#-running-with-docker)
- [API Endpoints](#-api-endpoints)
- [Configuration](#-configuration)
- [Testing](#-testing)
- [Example Test Addresses](#-example-test-addresses)
- [Development Notes](#-development-notes)
- [Troubleshooting](#-troubleshooting)

## 🏗️ Architecture

The project follows **Clean Architecture** principles, with a clear separation of concerns:

```
eth-tx-parser/
├── cmd/                  # Application entry points
├── internal/            
│   ├── domain/          # Core business logic and interfaces
│   ├── infrastructure/  # External implementations (Ethereum client, storage)
│   ├── application/     # Use cases and business rules
│   └── delivery/        # HTTP delivery mechanism
├── pkg/                 # Shared packages
└── test/               # Integration tests
```

### 🔑 Key Design Decisions

1. **Clean Architecture**
    - Clear dependency rules (dependencies point inward)
    - Core business logic is independent of external concerns
    - Easy to test and modify components

2. **In-Memory Storage**
    - Simple and fast for prototype/MVP
    - Easily replaceable with persistent storage
    - Thread-safe implementation

3. **HTTP API**
    - RESTful endpoints for easy integration
    - JSON responses
    - Simple subscription model

## 📋 Prerequisites

- Go 1.22 or compatible version
- Docker (optional, version 20.10.0 or higher recommended)
- Access to an Ethereum node (default uses public node)

## 🚀 Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ether-tx-parser
cd ether-tx-parser
```

2. Initialize the Go module:
```bash
go mod init ether-tx-parser
go mod tidy
```

3. Build and run locally:
```bash
go run cmd/api/main.go
```

## 🐳 Running with Docker

1. Build the Docker image:
```bash
docker build . -t ether-tx-parser
```

2. Run the container:
```bash
docker run -p 8080:8080 ether-tx-parser
```

Or use **docker-compose**:
```bash
docker-compose up
```

## 🌍 API Endpoints

### 1. Subscribe to an Address
```bash
curl -X POST http://localhost:8080/subscribe \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x28C6c06298d514Db089934071355E5743bf21d60"
  }'

# Expected Response:
# {"success":true}
```

### 2. Get Current Block
```bash
curl http://localhost:8080/block

# Expected Response:
# {"current_block":18934567}
```

### 3. Get Transactions
```bash
curl "http://localhost:8080/transactions?address=0x28C6c06298d514Db089934071355E5743bf21d60"

# Expected Response:
# [
#   {
#     "hash": "0x123...",
#     "from": "0x28C6c06298d514Db089934071355E5743bf21d60",
#     "to": "0x456...",
#     "value": "1000000000000000000",
#     "blockNumber": 18934566
#   }
# ]
```

## ⚙️ Configuration

The service can be configured through environment variables:

```bash
ETH_PARSER_SERVER_PORT=8080
ETH_PARSER_ETHEREUM_RPC_URL="https://ethereum-rpc.publicnode.com"
```

## 🧪 Testing

### Running Unit Tests
```bash
go test ./...
```

### Running Integration Tests

Locally:
1. Start the service:
```bash
go run cmd/api/main.go
```

2. In another terminal:
```bash
INTEGRATION_TEST=true go test ./test/integration -v
```

## 🔍 Example Test Addresses

For testing purposes, you can use these active Ethereum addresses:
- Binance Hot Wallet: `0x28C6c06298d514Db089934071355E5743bf21d60`
- Coinbase: `0x503828976D22510aad0201ac7EC88293211D23Da`

## 📌 Development Notes

- The service polls for new blocks every **15 seconds**
- Transactions are stored **in memory** and will be lost on service restart
- Address subscriptions are also stored in memory
- The service starts parsing from **10 blocks before** the current block on startup

## 🛠️ Troubleshooting

1. **If Docker build fails:**
    - Ensure `go.mod` and `go.sum` are present
    - Verify all required files are in the correct locations

2. **If the service isn't finding transactions:**
    - Verify the address is subscribed using the debug endpoint
    - Check the logs for any RPC errors
    - Ensure the Ethereum node URL is accessible

