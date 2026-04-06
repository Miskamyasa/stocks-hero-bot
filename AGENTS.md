# AGENTS.md

This document provides guidance for AI coding agents working in this repository.

## Project Overview

A Telegram bot that tracks stock portfolios and sends periodic balance updates. Built in Go with:
- Telegram Bot API (`github.com/go-telegram-bot-api/telegram-bot-api/v5`)
- Pure-Go SQLite (`modernc.org/sqlite`) - no CGO required
- Yahoo Finance API for stock data

## Project Structure

```
cmd/bot/main.go              # Entrypoint, config loading, dependency wiring
internal/
├── bot/                     # Telegram long-poll loop, FSM handlers
├── db/                      # SQLite connection, schema, CRUD operations
├── finance/                 # Yahoo Finance API client, price cache
├── portfolio/               # Business logic, balance computation
└── scheduler/               # Periodic notifications to users
```

## Build Commands

```bash
go run ./cmd/bot                                  # Run locally
go build -o stocks-hero-bot ./cmd/bot             # Build binary
make build                                        # Build for Linux (amd64)
go build ./...                                    # Verify compilation
```

## Lint Commands

```bash
go vet ./...                                      # Static analysis (required)
golangci-lint run                                 # Full linter suite
mise install                                      # Install tools via mise
```

## Test Commands

```bash
go test ./...                                     # Run all tests
go test -v ./...                                  # Verbose output
go test -v -run TestFunctionName ./path/to/pkg   # Run single test
go test ./internal/finance/...                    # Test specific package
go test -race ./...                               # Race detection
go test -cover ./...                              # Coverage report
```

## Code Style Guidelines

### Formatting

- Use `gofmt` (enforced by `golangci-lint`)
- Indentation: tabs for Go files, spaces for others (see `.editorconfig`)
- Line endings: LF
- Final newline: required

### Imports

Group imports in this order, separated by blank lines:
1. Standard library
2. External dependencies
3. Internal packages (prefixed with `stock-portfolio-bot/internal/`)

```go
import (
    "context"
    "fmt"
    "log"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

    "stock-portfolio-bot/internal/db"
    "stock-portfolio-bot/internal/finance"
)
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `bot`, `finance`, `scheduler`)
- **Exported types**: PascalCase with descriptive names (e.g., `YahooClient`, `BalanceReport`)
- **Unexported types**: camelCase (e.g., `yahooSession`, `cachedQuote`)
- **Interfaces**: describe behavior, often end in `-er` (e.g., `Notifier`)
- **Constructor functions**: `New<Type>` pattern (e.g., `NewYahooClient`, `NewService`)
- **Acronyms**: keep consistent casing (e.g., `chatID`, `TTL`, `URL`)

### Error Handling

- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Log errors at the handling site, not at every level
- Return `nil, nil` for "not found" cases (e.g., empty portfolio)
- Use early returns for error cases

```go
if err != nil {
    return nil, fmt.Errorf("get holdings: %w", err)
}

// Log at the handler level
if err != nil {
    log.Printf("compute balance %d: %v", chatID, err)
    h.sendText(chatID, "Failed to fetch prices. Please try again later.")
    return
}
```

### Concurrency

- Use `sync.Mutex` or `sync.RWMutex` for shared state
- Defer unlocks immediately after locking
- Use `context.Context` for cancellation and timeouts
- Launch goroutines with `go` keyword for concurrent operations

### Database

- Use parameterized queries (never string concatenation)
- Close rows with `defer func() { _ = rows.Close() }()`
- Handle `sql.ErrNoRows` explicitly for optional results
- Schema migrations are in `internal/db/sqlite.go`

### HTTP Clients

- Always set timeouts on `http.Client`
- Use `http.NewRequestWithContext` for cancellation support
- Drain and close response bodies
- Set appropriate User-Agent headers

### Telegram Bot Patterns

- FSM state stored in database (survives restarts)
- Callback data prefixed with action (e.g., `select:`, `remove:`)
- Acknowledge callbacks immediately to remove loading spinner
- Use Markdown formatting for rich messages

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | (required) | Token from @BotFather |
| `DB_PATH` | `./portfolio.db` | SQLite database path |
| `CACHE_TTL` | `15m` | Price cache duration |
| `NOTIFY_INTERVAL` | `30m` | Balance notification interval |

## Pre-commit Checklist

Before submitting changes:

```bash
go build ./...
go vet ./...
golangci-lint run
go test ./...
```
