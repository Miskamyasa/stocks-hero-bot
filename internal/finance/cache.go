package finance

import (
	"sync"
	"time"
)

// PriceCache is a thread-safe in-memory cache for stock quotes with TTL expiry.
// It is shared across all users so each symbol is fetched at most once per TTL window.
type PriceCache struct {
	mu    sync.RWMutex
	items map[string]cachedQuote
	ttl   time.Duration
}

type cachedQuote struct {
	quote     Quote
	fetchedAt time.Time
}

// ExchangeRateCache is a thread-safe in-memory cache for FX rates with TTL expiry.
// Rates are keyed by currency code and represent USD->currency conversions.
type ExchangeRateCache struct {
	mu    sync.RWMutex
	items map[string]cachedRate
	ttl   time.Duration
}

type cachedRate struct {
	rate      float64
	fetchedAt time.Time
}

// NewPriceCache creates a PriceCache with the given TTL.
func NewPriceCache(ttl time.Duration) *PriceCache {
	return &PriceCache{
		items: make(map[string]cachedQuote),
		ttl:   ttl,
	}
}

// Get returns a cached quote if it exists and has not expired.
func (pc *PriceCache) Get(symbol string) (Quote, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	item, ok := pc.items[symbol]
	if !ok || time.Since(item.fetchedAt) > pc.ttl {
		return Quote{}, false
	}
	return item.quote, true
}

// Set stores a quote with the current timestamp.
func (pc *PriceCache) Set(symbol string, q Quote) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.items[symbol] = cachedQuote{
		quote:     q,
		fetchedAt: time.Now(),
	}
}

// GetMulti performs a bulk cache lookup.
// Returns the cached quotes that are still fresh and the list of symbols that need fetching.
func (pc *PriceCache) GetMulti(symbols []string) (found map[string]Quote, missing []string) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	found = make(map[string]Quote)
	for _, sym := range symbols {
		item, ok := pc.items[sym]
		if ok && time.Since(item.fetchedAt) <= pc.ttl {
			found[sym] = item.quote
		} else {
			missing = append(missing, sym)
		}
	}
	return
}

// SetMulti stores multiple quotes at once.
func (pc *PriceCache) SetMulti(quotes map[string]Quote) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for sym, q := range quotes {
		pc.items[sym] = cachedQuote{quote: q, fetchedAt: now}
	}
}

// NewExchangeRateCache creates an ExchangeRateCache with the given TTL.
func NewExchangeRateCache(ttl time.Duration) *ExchangeRateCache {
	return &ExchangeRateCache{
		items: make(map[string]cachedRate),
		ttl:   ttl,
	}
}

// Get returns a cached rate if it exists and has not expired.
func (rc *ExchangeRateCache) Get(currency string) (float64, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	item, ok := rc.items[currency]
	if !ok || time.Since(item.fetchedAt) > rc.ttl {
		return 0, false
	}
	return item.rate, true
}

// Set stores a rate with the current timestamp.
func (rc *ExchangeRateCache) Set(currency string, rate float64) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.items[currency] = cachedRate{
		rate:      rate,
		fetchedAt: time.Now(),
	}
}

// GetMulti performs a bulk cache lookup.
// Returns the cached rates that are still fresh and the list of currency codes that need fetching.
func (rc *ExchangeRateCache) GetMulti(currencies []string) (found map[string]float64, missing []string) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	found = make(map[string]float64)
	for _, currency := range currencies {
		item, ok := rc.items[currency]
		if ok && time.Since(item.fetchedAt) <= rc.ttl {
			found[currency] = item.rate
		} else {
			missing = append(missing, currency)
		}
	}
	return
}

// SetMulti stores multiple rates at once.
func (rc *ExchangeRateCache) SetMulti(rates map[string]float64) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	now := time.Now()
	for currency, rate := range rates {
		rc.items[currency] = cachedRate{rate: rate, fetchedAt: now}
	}
}
