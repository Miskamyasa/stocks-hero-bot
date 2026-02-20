# stocks-hero-bot

A Telegram bot that tracks your stock portfolio and sends hourly balance updates.

Send a ticker or company name, pick from the search results, enter your share count — the bot handles the rest.

## Features

- Search stocks and ETFs by symbol or company name (Yahoo Finance)
- Track fractional shares across multiple positions
- Hourly portfolio balance notifications
- Per-user FSM conversation flow with persistent state (survives restarts)
- Global price cache with configurable TTL — one fetch per symbol regardless of user count
- Rate-limited fetch queue (no Yahoo API hammering)
- Pure-Go SQLite — no CGO, easy cross-compilation

## Bot commands

| Command | Description |
|---|---|
| _(any text)_ | Search for a ticker by symbol or company name |
| `/portfolio` | Show current holdings with live prices and total value |
| `/remove` | Remove a holding via inline buttons |
| `/start` | Show welcome message and reset state |
| `/help` | Show usage instructions |

## Requirements

- Go 1.22+
- A Telegram bot token from [@BotFather](https://t.me/BotFather)

No C compiler or system SQLite library needed — `modernc.org/sqlite` is pure Go.

## Running locally

```bash
git clone <repo>
cd stocks-hero-bot

cp .env.example .env
# edit .env and set TELEGRAM_BOT_TOKEN

go run ./cmd/bot
```

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | _(required)_ | Token from @BotFather |
| `DB_PATH` | `./portfolio.db` | Path to the SQLite database file |
| `CACHE_TTL` | `15m` | How long prices are cached before re-fetching |
| `RATE_LIMIT_PER_SEC` | `5` | Max Yahoo Finance requests per second |
| `NOTIFY_INTERVAL` | `1h` | How often to push balance updates to users |

### Building a binary

```bash
go build -o stocks-hero-bot ./cmd/bot
./stocks-hero-bot
```

Cross-compile for Linux (e.g. for a VPS):

```bash
GOOS=linux GOARCH=amd64 go build -o stocks-hero-bot-linux ./cmd/bot
```

## Dev container

The repo includes a `.devcontainer/devcontainer.json` for VS Code / GitHub Codespaces. Open the folder in VS Code and choose **Reopen in Container** — `mise install` runs automatically and installs Go and `golangci-lint`.

## Project structure

```
stocks-hero-bot/
├── cmd/bot/
│   └── main.go              # entrypoint, config loading
├── internal/
│   ├── bot/
│   │   ├── bot.go           # Telegram long-poll loop, SendMarkdown
│   │   └── handler.go       # FSM message and callback handlers
│   ├── db/
│   │   ├── sqlite.go        # connection, schema migration
│   │   └── repository.go    # CRUD: users, holdings
│   ├── finance/
│   │   ├── yahoo.go         # search and batch quote endpoints
│   │   ├── cache.go         # TTL price cache shared across all users
│   │   ├── queue.go         # rate-limited fetch queue with batching
│   │   └── http.go          # shared http.Client
│   ├── portfolio/
│   │   └── service.go       # ComputeBalance, BalanceReport formatting
│   └── scheduler/
│       └── scheduler.go     # hourly tick → pre-warm cache → notify users
├── .env.example
├── mise.toml                # tool versions (Go, golangci-lint)
└── go.mod
```

## Contributing

1. Fork the repo and create a branch from `main`.
2. Install tools with [mise](https://mise.jdx.dev/): `mise install`
3. Make your changes. Keep each commit focused on one thing.
4. Run checks before opening a PR:

```bash
go build ./...
go vet ./...
golangci-lint run
```

5. Open a pull request with a clear description of what changed and why.

### Code style

- Standard `gofmt` formatting (enforced by `golangci-lint`).
- Errors are wrapped with `fmt.Errorf("context: %w", err)` and logged at the call site that handles them — not re-logged up the stack.
- New behaviour should be covered by a test where practical.

## License

GPL-3.0 — see [LICENSE](LICENSE).
