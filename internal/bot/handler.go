package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"stock-portfolio-bot/internal/db"
	"stock-portfolio-bot/internal/finance"
	"stock-portfolio-bot/internal/portfolio"
)

const welcomeText = `Welcome! I'm your stock portfolio assistant 📈

Here's what I can do:
• Send me a ticker symbol or company name — I'll look it up
• Select the right match from the list
• Tell me how many shares you own
• I'll track prices and notify you every hour with your total balance

Let's start — send me a ticker symbol or company name!`

// Handler processes Telegram messages and callbacks using a per-user FSM.
type Handler struct {
	api   *tgbotapi.BotAPI
	svc   *portfolio.Service
	yahoo *finance.YahooClient
	repo  *db.Repository
}

func newHandler(api *tgbotapi.BotAPI, svc *portfolio.Service, yahoo *finance.YahooClient) *Handler {
	return &Handler{
		api:   api,
		svc:   svc,
		yahoo: yahoo,
		repo:  svc.Repo(),
	}
}

// HandleMessage routes an incoming text message based on the user's FSM state.
func (h *Handler) HandleMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	if err := h.repo.UpsertUser(chatID, msg.From.UserName); err != nil {
		log.Printf("upsert user %d: %v", chatID, err)
	}

	// Handle commands first.
	if msg.IsCommand() {
		h.handleCommand(ctx, msg)
		return
	}

	state, stateData, err := h.repo.GetUserState(chatID)
	if err != nil {
		log.Printf("get user state %d: %v", chatID, err)
		return
	}

	switch state {
	case "idle", "", "awaiting_ticker_choice":
		h.handleTickerSearch(ctx, chatID, msg.Text)

	case "awaiting_shares":
		h.handleSharesInput(ctx, chatID, msg.Text, stateData)

	default:
		h.sendText(chatID, "Please select a ticker from the list above, or type a new ticker to search.")
	}
}

// HandleCallback routes an inline keyboard callback.
func (h *Handler) HandleCallback(ctx context.Context, cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	data := cb.Data

	// Acknowledge the callback to remove the loading spinner.
	ack := tgbotapi.NewCallback(cb.ID, "")
	if _, err := h.api.Request(ack); err != nil {
		log.Printf("ack callback: %v", err)
	}

	switch {
	case strings.HasPrefix(data, "select:"):
		symbol := strings.TrimPrefix(data, "select:")
		h.handleTickerSelect(ctx, chatID, symbol)

	case strings.HasPrefix(data, "remove:"):
		symbol := strings.TrimPrefix(data, "remove:")
		h.handleRemove(ctx, chatID, symbol)
	}
}

// --- command handlers ---

func (h *Handler) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	switch msg.Command() {
	case "start":
		_ = h.repo.SetUserState(chatID, "idle", "")
		h.sendText(chatID, welcomeText)

	case "portfolio":
		h.handlePortfolio(ctx, chatID)

	case "remove":
		h.handleRemoveMenu(ctx, chatID)

	case "help":
		h.sendText(chatID, welcomeText)

	default:
		h.sendText(chatID, "Unknown command. Use /portfolio, /remove, or /help.")
	}
}

func (h *Handler) handlePortfolio(ctx context.Context, chatID int64) {
	report, err := h.svc.ComputeBalance(ctx, chatID)
	if err != nil {
		log.Printf("compute balance %d: %v", chatID, err)
		h.sendText(chatID, "Failed to fetch prices. Please try again later.")
		return
	}
	if report == nil || len(report.Holdings) == 0 {
		h.sendText(chatID, "Your portfolio is empty. Send me a ticker symbol to get started!")
		return
	}
	h.sendMarkdown(chatID, report.Format())
}

func (h *Handler) handleRemoveMenu(ctx context.Context, chatID int64) {
	holdings, err := h.repo.GetHoldings(chatID)
	if err != nil || len(holdings) == 0 {
		h.sendText(chatID, "Your portfolio is empty.")
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, holding := range holdings {
		label := fmt.Sprintf("❌ %s — %s", holding.Symbol, holding.Name)
		btn := tgbotapi.NewInlineKeyboardButtonData(label, "remove:"+holding.Symbol)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewMessage(chatID, "Select a holding to remove:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send remove menu %d: %v", chatID, err)
	}
}

// --- FSM state handlers ---

func (h *Handler) handleTickerSearch(ctx context.Context, chatID int64, query string) {
	if strings.TrimSpace(query) == "" {
		h.sendText(chatID, "Please send a ticker symbol or company name.")
		return
	}

	results, err := h.yahoo.SearchTickers(ctx, query)
	if err != nil {
		log.Printf("search tickers %q: %v", query, err)
		h.sendText(chatID, "Search failed. Please try again.")
		return
	}
	if len(results) == 0 {
		h.sendText(chatID, "No tickers found. Try another name or symbol.")
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, r := range results {
		label := fmt.Sprintf("%s — %s (%s)", r.Symbol, r.Name, r.Exchange)
		btn := tgbotapi.NewInlineKeyboardButtonData(label, "select:"+r.Symbol)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	msg := tgbotapi.NewMessage(chatID, "Select a ticker:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send ticker list %d: %v", chatID, err)
		return
	}

	stateJSON, _ := json.Marshal(results)
	if err := h.repo.SetUserState(chatID, "awaiting_ticker_choice", string(stateJSON)); err != nil {
		log.Printf("set user state %d: %v", chatID, err)
	}
}

func (h *Handler) handleTickerSelect(ctx context.Context, chatID int64, symbol string) {
	_, stateData, err := h.repo.GetUserState(chatID)
	if err != nil {
		log.Printf("get user state %d: %v", chatID, err)
		return
	}

	// Recover the name from the stored search results.
	var results []finance.TickerResult
	name := symbol // fallback
	if err := json.Unmarshal([]byte(stateData), &results); err == nil {
		for _, r := range results {
			if r.Symbol == symbol {
				name = r.Name
				break
			}
		}
	}

	pending := struct {
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}{Symbol: symbol, Name: name}

	pendingJSON, _ := json.Marshal(pending)
	if err := h.repo.SetUserState(chatID, "awaiting_shares", string(pendingJSON)); err != nil {
		log.Printf("set user state %d: %v", chatID, err)
		return
	}

	h.sendText(chatID, fmt.Sprintf(
		"You selected *%s* (%s).\n\nHow many shares do you own? (fractional shares are supported)",
		symbol, name,
	))
}

func (h *Handler) handleSharesInput(ctx context.Context, chatID int64, text, stateData string) {
	shares, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil || shares <= 0 {
		h.sendText(chatID, "Please enter a valid positive number of shares (e.g. 10 or 2.5).")
		return
	}

	var pending struct {
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}
	if err := json.Unmarshal([]byte(stateData), &pending); err != nil {
		log.Printf("unmarshal pending state %d: %v", chatID, err)
		h.sendText(chatID, "Something went wrong. Please start over by sending a ticker symbol.")
		_ = h.repo.SetUserState(chatID, "idle", "")
		return
	}

	if err := h.repo.UpsertHolding(chatID, pending.Symbol, pending.Name, shares); err != nil {
		log.Printf("upsert holding %d %s: %v", chatID, pending.Symbol, err)
		h.sendText(chatID, "Failed to save holding. Please try again.")
		return
	}

	if err := h.repo.SetUserState(chatID, "idle", ""); err != nil {
		log.Printf("reset user state %d: %v", chatID, err)
	}

	h.sendText(chatID, fmt.Sprintf(
		"✅ Saved: %.4f shares of %s (%s).\n\nSend another ticker to add more, or /portfolio to see your balance.",
		shares, pending.Symbol, pending.Name,
	))
}

func (h *Handler) handleRemove(ctx context.Context, chatID int64, symbol string) {
	if err := h.repo.DeleteHolding(chatID, symbol); err != nil {
		log.Printf("delete holding %d %s: %v", chatID, symbol, err)
		h.sendText(chatID, "Failed to remove holding. Please try again.")
		return
	}
	h.sendText(chatID, fmt.Sprintf("✅ Removed %s from your portfolio.", symbol))
}

// --- helpers ---

func (h *Handler) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send text %d: %v", chatID, err)
	}
}

func (h *Handler) sendMarkdown(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := h.api.Send(msg); err != nil {
		log.Printf("send markdown %d: %v", chatID, err)
	}
}
