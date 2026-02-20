package finance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	searchURL  = "https://query2.finance.yahoo.com/v1/finance/search"
	chartURL   = "https://query1.finance.yahoo.com/v8/finance/chart"
	consentURL = "https://fc.yahoo.com/"
	crumbURL   = "https://query2.finance.yahoo.com/v1/test/getcrumb"

	sessionTTL = 30 * time.Minute
)

// TickerResult is a single search result from Yahoo Finance.
type TickerResult struct {
	Symbol   string
	Name     string
	Exchange string
	Type     string // "EQUITY", "ETF", etc.
}

// Quote holds the latest price data for a symbol.
type Quote struct {
	Symbol   string
	Price    float64
	Currency string
}

// yahooSession holds the cookie and crumb required by Yahoo Finance API.
type yahooSession struct {
	cookie    string
	crumb     string
	expiresAt time.Time
}

// YahooClient fetches data from Yahoo Finance with session-based auth and a shared price cache.
type YahooClient struct {
	cache  *PriceCache
	client *http.Client

	sessionMu sync.Mutex
	session   *yahooSession
}

// NewYahooClient creates a YahooClient backed by the given PriceCache.
func NewYahooClient(cache *PriceCache) *YahooClient {
	return &YahooClient{
		cache:  cache,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// --- session management ---

func (yc *YahooClient) getSession(ctx context.Context) (*yahooSession, error) {
	yc.sessionMu.Lock()
	defer yc.sessionMu.Unlock()

	if yc.session != nil && time.Now().Before(yc.session.expiresAt) {
		return yc.session, nil
	}

	sess, err := yc.fetchNewSession(ctx)
	if err != nil {
		return nil, err
	}
	yc.session = sess
	return sess, nil
}

func (yc *YahooClient) invalidateSession() {
	yc.sessionMu.Lock()
	yc.session = nil
	yc.sessionMu.Unlock()
}

func (yc *YahooClient) fetchNewSession(ctx context.Context) (*yahooSession, error) {
	// Step 1: hit fc.yahoo.com to get a consent cookie.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, consentURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build consent request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// We need the raw Set-Cookie headers; don't follow redirects automatically.
	noRedirectClient := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	consentResp, err := noRedirectClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("consent request: %w", err)
	}
	defer consentResp.Body.Close()
	_, _ = io.ReadAll(consentResp.Body) // drain

	cookie := extractCookies(consentResp)
	if cookie == "" {
		return nil, fmt.Errorf("no cookie returned from Yahoo consent endpoint")
	}

	// Step 2: exchange the cookie for a crumb.
	crumbReq, err := http.NewRequestWithContext(ctx, http.MethodGet, crumbURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build crumb request: %w", err)
	}
	crumbReq.Header.Set("Cookie", cookie)
	crumbReq.Header.Set("User-Agent", "Mozilla/5.0")

	crumbResp, err := yc.client.Do(crumbReq)
	if err != nil {
		return nil, fmt.Errorf("crumb request: %w", err)
	}
	defer crumbResp.Body.Close()

	crumbBytes, err := io.ReadAll(crumbResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read crumb response: %w", err)
	}

	crumb := strings.TrimSpace(string(crumbBytes))
	if crumb == "" || crumb == "null" {
		return nil, fmt.Errorf("Yahoo returned invalid crumb: %q", crumb)
	}

	return &yahooSession{
		cookie:    cookie,
		crumb:     crumb,
		expiresAt: time.Now().Add(sessionTTL),
	}, nil
}

// extractCookies collects all Set-Cookie name=value pairs into a single Cookie header string.
func extractCookies(resp *http.Response) string {
	var parts []string
	for _, c := range resp.Cookies() {
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}

// --- public API ---

// SearchTickers queries Yahoo Finance for tickers matching query.
// Results are filtered to EQUITY and ETF types only.
// Search does not require authentication.
func (yc *YahooClient) SearchTickers(ctx context.Context, query string) ([]TickerResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("quotesCount", "8")
	params.Set("newsCount", "0")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("build search request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := yc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read search response: %w", err)
	}

	var payload struct {
		Quotes []struct {
			Symbol    string `json:"symbol"`
			Shortname string `json:"shortname"`
			Longname  string `json:"longname"`
			Exchange  string `json:"exchange"`
			QuoteType string `json:"quoteType"`
		} `json:"quotes"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	var results []TickerResult
	for _, q := range payload.Quotes {
		if q.QuoteType != "EQUITY" && q.QuoteType != "ETF" {
			continue
		}
		name := q.Shortname
		if name == "" {
			name = q.Longname
		}
		results = append(results, TickerResult{
			Symbol:   q.Symbol,
			Name:     name,
			Exchange: q.Exchange,
			Type:     q.QuoteType,
		})
	}
	return results, nil
}

// GetQuotes returns prices for the given symbols, using the cache where fresh.
func (yc *YahooClient) GetQuotes(ctx context.Context, symbols []string) (map[string]Quote, error) {
	if len(symbols) == 0 {
		return map[string]Quote{}, nil
	}

	found, missing := yc.cache.GetMulti(symbols)
	if len(missing) == 0 {
		return found, nil
	}

	fetched, err := yc.fetchBatch(ctx, missing)
	if err != nil {
		return found, err
	}

	yc.cache.SetMulti(fetched)
	for k, v := range fetched {
		found[k] = v
	}
	return found, nil
}

// fetchBatch fetches prices for multiple symbols using the v8/chart endpoint,
// one request per symbol (the chart endpoint is per-symbol, not batch).
// It retries once with a fresh session on 401/403.
func (yc *YahooClient) fetchBatch(ctx context.Context, symbols []string) (map[string]Quote, error) {
	sess, err := yc.getSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("get yahoo session: %w", err)
	}

	quotes := make(map[string]Quote, len(symbols))
	for _, sym := range symbols {
		q, err := yc.fetchOne(ctx, sym, sess)
		if err != nil {
			// On auth failure, refresh session and retry once.
			if isAuthError(err) {
				yc.invalidateSession()
				sess, err = yc.getSession(ctx)
				if err != nil {
					return quotes, fmt.Errorf("refresh yahoo session: %w", err)
				}
				q, err = yc.fetchOne(ctx, sym, sess)
			}
			if err != nil {
				return quotes, fmt.Errorf("fetch %s: %w", sym, err)
			}
		}
		quotes[sym] = q
	}
	return quotes, nil
}

func (yc *YahooClient) fetchOne(ctx context.Context, symbol string, sess *yahooSession) (Quote, error) {
	u := fmt.Sprintf("%s/%s?range=1d&interval=1d&crumb=%s",
		chartURL, url.PathEscape(symbol), url.QueryEscape(sess.crumb))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Quote{}, fmt.Errorf("build chart request: %w", err)
	}
	req.Header.Set("Cookie", sess.cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := yc.client.Do(req)
	if err != nil {
		return Quote{}, fmt.Errorf("chart request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return Quote{}, fmt.Errorf("auth error: HTTP %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return Quote{}, fmt.Errorf("HTTP %d for %s", resp.StatusCode, symbol)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Quote{}, fmt.Errorf("read chart response: %w", err)
	}

	return parseChartResponse(body, symbol)
}

func parseChartResponse(body []byte, symbol string) (Quote, error) {
	var payload struct {
		Chart struct {
			Result []struct {
				Meta struct {
					Symbol             string  `json:"symbol"`
					Currency           string  `json:"currency"`
					RegularMarketPrice float64 `json:"regularMarketPrice"`
					ChartPreviousClose float64 `json:"chartPreviousClose"`
				} `json:"meta"`
			} `json:"result"`
			Error *struct {
				Description string `json:"description"`
			} `json:"error"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return Quote{}, fmt.Errorf("parse chart response: %w", err)
	}

	if payload.Chart.Error != nil {
		return Quote{}, fmt.Errorf("yahoo error: %s", payload.Chart.Error.Description)
	}

	if len(payload.Chart.Result) == 0 {
		return Quote{}, fmt.Errorf("no chart result for %s", symbol)
	}

	meta := payload.Chart.Result[0].Meta
	price := meta.RegularMarketPrice
	if price == 0 {
		price = meta.ChartPreviousClose // fallback for closed markets
	}
	if price == 0 {
		return Quote{}, fmt.Errorf("no price data for %s", symbol)
	}

	return Quote{
		Symbol:   meta.Symbol,
		Price:    price,
		Currency: meta.Currency,
	}, nil
}

func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "auth error") ||
		strings.Contains(s, "HTTP 401") ||
		strings.Contains(s, "HTTP 403")
}
