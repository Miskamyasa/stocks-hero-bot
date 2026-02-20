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
