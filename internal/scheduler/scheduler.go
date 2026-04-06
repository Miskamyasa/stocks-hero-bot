package scheduler

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"stock-portfolio-bot/internal/portfolio"
)

// changeThreshold is the minimum percentage change (as a decimal) required
// to trigger a notification. 0.0005 = 0.05%.
const changeThreshold = 0.0005

// Notifier is the interface the scheduler uses to push messages to users.
// Implemented by *bot.Bot to avoid an import cycle.
type Notifier interface {
	SendMarkdown(chatID int64, text string)
}

// Scheduler fires periodic portfolio notifications for all active users.
type Scheduler struct {
	svc      *portfolio.Service
	notifier Notifier
	interval time.Duration
}

// New creates a Scheduler.
func New(svc *portfolio.Service, notifier Notifier, interval time.Duration) *Scheduler {
	return &Scheduler{svc: svc, notifier: notifier, interval: interval}
}

// Run starts the notification loop. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.notifyAll(ctx)
		}
	}
}

func (s *Scheduler) notifyAll(ctx context.Context) {
	repo := s.svc.Repo()

	// 1. Pre-warm cache: batch-fetch all distinct symbols once.
	symbols, err := repo.GetDistinctSymbols()
	if err != nil {
		log.Printf("scheduler: get distinct symbols: %v", err)
		return
	}
	if len(symbols) == 0 {
		return
	}

	quotes, err := s.svc.GetQuotes(ctx, symbols)
	if err != nil {
		log.Printf("scheduler: pre-warm quotes: %v", err)
		// Continue anyway — individual balance calls will retry.
	}

	// 2. Pre-warm FX cache for currencies from the current quote universe.
	currencySet := make(map[string]struct{})
	for _, q := range quotes {
		currency := strings.ToUpper(strings.TrimSpace(q.Currency))
		if currency == "" || currency == "USD" {
			continue
		}
		currencySet[currency] = struct{}{}
	}

	if len(currencySet) > 0 {
		currencies := make([]string, 0, len(currencySet))
		for currency := range currencySet {
			currencies = append(currencies, currency)
		}
		if err := s.svc.PrewarmUSDRates(ctx, currencies); err != nil {
			log.Printf("scheduler: pre-warm USD rates: %v", err)
			// Continue anyway — per-user computations will retry.
		}
	}

	// 3. Notify each active user (balance reads from cache → instant).
	users, err := repo.GetAllActiveUsers()
	if err != nil {
		log.Printf("scheduler: get active users: %v", err)
		return
	}

	for _, chatID := range users {
		report, err := s.svc.ComputeBalance(ctx, chatID)
		if err != nil {
			log.Printf("scheduler: compute balance %d: %v", chatID, err)
			continue
		}
		if report == nil || len(report.Holdings) == 0 {
			continue
		}

		// Fetch previous report to detect changes.
		prev, err := repo.GetLastReport(chatID)
		if err != nil {
			log.Printf("scheduler: get last report %d: %v", chatID, err)
		}

		// Skip sending and saving if the change is below threshold.
		if prev > 0 {
			changeRatio := math.Abs(report.TotalUSD-prev) / prev
			if changeRatio < changeThreshold {
				log.Printf("scheduler: skip report for %d (change %.4f%% < threshold %.2f%%)",
					chatID, changeRatio*100, changeThreshold*100)
				continue
			}
		}

		text := report.FormatSummary()

		// Append % change vs previous report if one exists.
		if prev > 0 {
			change := (report.TotalUSD - prev) / prev * 100
			sign := "+"
			if change < 0 {
				sign = ""
			}
			text += fmt.Sprintf("\n📈 *Change since last report: %s%.2f%%*", sign, change)
		}

		s.notifier.SendMarkdown(chatID, text)

		if err := repo.SaveReport(chatID, report.TotalUSD); err != nil {
			log.Printf("scheduler: save report %d: %v", chatID, err)
		}
	}
}
