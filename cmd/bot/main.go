package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"stock-portfolio-bot/internal/bot"
	"stock-portfolio-bot/internal/db"
	"stock-portfolio-bot/internal/finance"
	"stock-portfolio-bot/internal/portfolio"
	"stock-portfolio-bot/internal/scheduler"
)

type config struct {
	TelegramToken  string
	DBPath         string
	CacheTTL       time.Duration
	NotifyInterval time.Duration
}

func loadConfig() config {
	_ = godotenv.Load() // ignore error if .env absent

	cacheTTL, err := time.ParseDuration(getEnv("CACHE_TTL", "15m"))
	if err != nil {
		cacheTTL = 15 * time.Minute
	}

	notifyInterval, err := time.ParseDuration(getEnv("NOTIFY_INTERVAL", "1h"))
	if err != nil {
		notifyInterval = time.Hour
	}

	return config{
		TelegramToken:  mustEnv("TELEGRAM_BOT_TOKEN"),
		DBPath:         getEnv("DB_PATH", "./portfolio.db"),
		CacheTTL:       cacheTTL,
		NotifyInterval: notifyInterval,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func main() {
	cfg := loadConfig()

	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer func() { _ = database.Close() }()

	priceCache := finance.NewPriceCache(cfg.CacheTTL)
	yahooClient := finance.NewYahooClient(priceCache)

	repo := db.NewRepository(database)
	svc := portfolio.NewService(repo, yahooClient)

	tgBot, err := bot.New(cfg.TelegramToken, svc, yahooClient)
	if err != nil {
		log.Fatalf("bot init: %v", err)
	}

	sched := scheduler.New(svc, tgBot, cfg.NotifyInterval)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go sched.Run(ctx)
	tgBot.Start(ctx)
}
