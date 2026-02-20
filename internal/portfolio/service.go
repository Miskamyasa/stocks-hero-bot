package portfolio

import (
	"context"
	"fmt"
	"log"
	"strings"

	"stock-portfolio-bot/internal/db"
	"stock-portfolio-bot/internal/finance"
)

// HoldingLine is one row in a balance report.
type HoldingLine struct {
	Symbol string
	Name   string
	Shares float64
	Price  float64
	Value  float64
}

// BalanceReport is the computed portfolio snapshot for a user.
type BalanceReport struct {
	Holdings []HoldingLine
	TotalUSD float64
}

// Format produces a Telegram-friendly Markdown message.
func (r *BalanceReport) Format() string {
	var sb strings.Builder
	sb.WriteString("📊 *Portfolio Balance*\n\n")
	for _, h := range r.Holdings {
		pct := 0.0
		if r.TotalUSD > 0 {
			pct = h.Value / r.TotalUSD * 100
		}
		sb.WriteString(fmt.Sprintf(
			"*%s* (%s)\n  %.4f shares × $%.2f = *$%.2f* (%.1f%%)\n",
			h.Symbol, h.Name, h.Shares, h.Price, h.Value, pct,
		))
	}
	sb.WriteString(fmt.Sprintf("\n💰 *Total: $%.2f*", r.TotalUSD))
	return sb.String()
}

// Service implements portfolio business logic.
type Service struct {
	repo  *db.Repository
	yahoo *finance.YahooClient
}

// NewService creates a Service.
func NewService(repo *db.Repository, yahoo *finance.YahooClient) *Service {
	return &Service{repo: repo, yahoo: yahoo}
}

// Repo exposes the repository (used by the scheduler).
func (s *Service) Repo() *db.Repository { return s.repo }

// GetQuotes delegates to the Yahoo client (used by the scheduler for cache pre-warming).
func (s *Service) GetQuotes(ctx context.Context, symbols []string) (map[string]finance.Quote, error) {
	return s.yahoo.GetQuotes(ctx, symbols)
}

// ComputeBalance fetches the latest prices and computes the total portfolio value for a user.
func (s *Service) ComputeBalance(ctx context.Context, chatID int64) (*BalanceReport, error) {
	holdings, err := s.repo.GetHoldings(chatID)
	if err != nil {
		return nil, fmt.Errorf("get holdings: %w", err)
	}
	if len(holdings) == 0 {
		return nil, nil
	}

	symbols := make([]string, len(holdings))
	for i, h := range holdings {
		symbols[i] = h.Symbol
	}

	quotes, err := s.yahoo.GetQuotes(ctx, symbols)
	if err != nil {
		return nil, fmt.Errorf("get quotes: %w", err)
	}

	report := &BalanceReport{Holdings: make([]HoldingLine, 0, len(holdings))}
	for _, h := range holdings {
		q, ok := quotes[h.Symbol]
		if !ok {
			log.Printf("ComputeBalance: no quote returned for symbol %s (chatID %d)", h.Symbol, chatID)
			continue
		}
		value := h.Shares * q.Price
		report.Holdings = append(report.Holdings, HoldingLine{
			Symbol: h.Symbol,
			Name:   h.Name,
			Shares: h.Shares,
			Price:  q.Price,
			Value:  value,
		})
		report.TotalUSD += value
	}
	if len(report.Holdings) == 0 {
		return nil, fmt.Errorf("quotes returned no data for symbols %v", symbols)
	}
	return report, nil
}
