package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"stock-portfolio-bot/internal/portfolio"
)

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

	if _, err := s.svc.GetQuotes(ctx, symbols); err != nil {
		log.Printf("scheduler: pre-warm quotes: %v", err)
		// Continue anyway — individual balance calls will retry.
	}

	// 2. Notify each active user (balance reads from cache → instant).
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

		text := report.Format()

		// Append % change vs previous report if one exists.
		prev, err := repo.GetLastReport(chatID)
		if err != nil {
			log.Printf("scheduler: get last report %d: %v", chatID, err)
		} else if prev > 0 {
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
