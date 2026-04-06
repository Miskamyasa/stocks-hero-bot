package portfolio

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"stock-portfolio-bot/internal/db"
	"stock-portfolio-bot/internal/finance"
)

// HoldingLine is one row in a balance report.
type HoldingLine struct {
	Symbol   string
	Name     string
	Shares   float64
	Price    float64
	Currency string
	Value    float64
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
		currency := strings.ToUpper(strings.TrimSpace(h.Currency))
		if currency == "" || currency == "USD" {
			fmt.Fprintf(&sb,
				"*%s* (%s)\n  %.4f shares × $%.2f = *$%.2f* (%.1f%%)\n",
				h.Symbol, h.Name, h.Shares, h.Price, h.Value, pct,
			)
			continue
		}

		fmt.Fprintf(&sb,
			"*%s* (%s)\n  %.4f shares × %.2f %s (%s->USD) = *$%.2f* (%.1f%%)\n",
			h.Symbol, h.Name, h.Shares, h.Price, currency, currency, h.Value, pct,
		)
	}
	fmt.Fprintf(&sb, "\n💰 *Total: $%.2f*", r.TotalUSD)
	return sb.String()
}

// FormatSummary returns only the total balance line in Markdown format.
func (r *BalanceReport) FormatSummary() string {
	return fmt.Sprintf("💰 *Total: $%.2f*", r.TotalUSD)
}

// Service implements portfolio business logic.
type Service struct {
	repo  *db.Repository
	yahoo *finance.YahooClient
	rates *finance.ExchangeRateCache
}

// NewService creates a Service.
func NewService(repo *db.Repository, yahoo *finance.YahooClient, rates *finance.ExchangeRateCache) *Service {
	return &Service{repo: repo, yahoo: yahoo, rates: rates}
}

// Repo exposes the repository (used by the scheduler).
func (s *Service) Repo() *db.Repository { return s.repo }

// GetQuotes delegates to the Yahoo client (used by the scheduler for cache pre-warming).
func (s *Service) GetQuotes(ctx context.Context, symbols []string) (map[string]finance.Quote, error) {
	return s.yahoo.GetQuotes(ctx, symbols)
}

// PrewarmUSDRates fetches and caches missing USD exchange rates for currencies.
func (s *Service) PrewarmUSDRates(ctx context.Context, currencies []string) error {
	if len(currencies) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(currencies))
	seen := make(map[string]struct{}, len(currencies))
	for _, code := range currencies {
		currency := strings.ToUpper(strings.TrimSpace(code))
		if currency == "" || currency == "USD" {
			continue
		}
		if _, ok := seen[currency]; ok {
			continue
		}
		seen[currency] = struct{}{}
		normalized = append(normalized, currency)
	}

	if len(normalized) == 0 {
		return nil
	}

	missing := normalized
	if s.rates != nil {
		_, missing = s.rates.GetMulti(normalized)
	}

	if len(missing) == 0 {
		return nil
	}

	fetchedRates, err := s.yahoo.GetUSDRates(ctx, missing)
	if s.rates != nil && len(fetchedRates) > 0 {
		s.rates.SetMulti(fetchedRates)
	}
	if err != nil {
		return fmt.Errorf("get USD rates: %w", err)
	}

	return nil
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

	currencySet := make(map[string]struct{})
	for _, h := range holdings {
		q, ok := quotes[h.Symbol]
		if !ok {
			continue
		}
		currency := strings.ToUpper(strings.TrimSpace(q.Currency))
		if currency == "" || currency == "USD" {
			continue
		}
		currencySet[currency] = struct{}{}
	}

	usdRates := map[string]float64{"USD": 1}
	if len(currencySet) > 0 {
		currencies := make([]string, 0, len(currencySet))
		for currency := range currencySet {
			currencies = append(currencies, currency)
		}

		missing := currencies
		if s.rates != nil {
			cachedRates, missingRates := s.rates.GetMulti(currencies)
			for currency, rate := range cachedRates {
				usdRates[currency] = rate
			}
			missing = missingRates
		}

		if len(missing) > 0 {
			fetchedRates, rateErr := s.yahoo.GetUSDRates(ctx, missing)
			for currency, rate := range fetchedRates {
				usdRates[currency] = rate
			}
			if s.rates != nil && len(fetchedRates) > 0 {
				s.rates.SetMulti(fetchedRates)
			}
			if rateErr != nil {
				log.Printf("ComputeBalance: get USD rates for currencies %v (chatID %d): %v", missing, chatID, rateErr)
			}
		}
	}

	report := &BalanceReport{Holdings: make([]HoldingLine, 0, len(holdings))}
	missingConversionCurrencies := make(map[string]struct{})
	for _, h := range holdings {
		q, ok := quotes[h.Symbol]
		if !ok {
			log.Printf("ComputeBalance: no quote returned for symbol %s (chatID %d)", h.Symbol, chatID)
			continue
		}
		valueUSD, ok := s.ConvertToUSD(h.Shares*q.Price, q.Currency, usdRates)
		if !ok {
			currency := strings.ToUpper(strings.TrimSpace(q.Currency))
			if currency == "" {
				currency = "UNKNOWN"
			}
			missingConversionCurrencies[currency] = struct{}{}
			log.Printf("ComputeBalance: no USD conversion rate for currency %s symbol %s (chatID %d)", currency, h.Symbol, chatID)
			continue
		}
		report.Holdings = append(report.Holdings, HoldingLine{
			Symbol:   h.Symbol,
			Name:     h.Name,
			Shares:   h.Shares,
			Price:    q.Price,
			Currency: q.Currency,
			Value:    valueUSD,
		})
		report.TotalUSD += valueUSD
	}
	if len(missingConversionCurrencies) > 0 {
		currencies := make([]string, 0, len(missingConversionCurrencies))
		for currency := range missingConversionCurrencies {
			currencies = append(currencies, currency)
		}
		sort.Strings(currencies)
		return nil, fmt.Errorf("missing USD conversion rates for currencies %v", currencies)
	}
	if len(report.Holdings) == 0 {
		return nil, fmt.Errorf("quotes returned no data for symbols %v", symbols)
	}
	return report, nil
}

// ResetBaseline computes the current portfolio balance, saves it as a new baseline
// for future performance comparisons, and returns the report along with the previous
// baseline value. This should be called after portfolio composition changes (adding/removing stocks)
// to ensure subsequent scheduler reports only reflect true performance changes.
// Returns (nil, 0, nil) if the portfolio is empty.
func (s *Service) ResetBaseline(ctx context.Context, chatID int64) (*BalanceReport, float64, error) {
	report, err := s.ComputeBalance(ctx, chatID)
	if err != nil {
		return nil, 0, fmt.Errorf("compute balance: %w", err)
	}
	if report == nil || len(report.Holdings) == 0 {
		return nil, 0, nil
	}

	prevTotal, err := s.repo.GetLastReport(chatID)
	if err != nil {
		return nil, 0, fmt.Errorf("get last report: %w", err)
	}

	if err := s.repo.SaveReport(chatID, report.TotalUSD); err != nil {
		return nil, 0, fmt.Errorf("save report: %w", err)
	}

	return report, prevTotal, nil
}

// ConvertToUSD converts a value from the quote currency to USD.
// usdRates map values are USD->currency conversions.
func (s *Service) ConvertToUSD(value float64, currency string, usdRates map[string]float64) (float64, bool) {
	normalizedCurrency := strings.ToUpper(strings.TrimSpace(currency))
	if normalizedCurrency == "USD" {
		return value, true
	}
	if normalizedCurrency == "" {
		return 0, false
	}

	rate, ok := usdRates[normalizedCurrency]
	if !ok || rate <= 0 {
		return 0, false
	}

	return value / rate, true
}
