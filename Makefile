.PHONY: build

## build: compile the bot binary for Linux amd64 to /tmp/stocks-hero-bot
build:
	GOOS=linux GOARCH=amd64 go build -o stocks-hero-bot ./cmd/bot
	chmod +x stocks-hero-bot
