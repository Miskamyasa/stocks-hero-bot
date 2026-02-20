package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"stock-portfolio-bot/internal/finance"
	"stock-portfolio-bot/internal/portfolio"
)

// Bot wraps the Telegram bot API and owns the update dispatch loop.
type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
}

// botCommands is the list registered with Telegram so they appear in the menu.
var botCommands = []tgbotapi.BotCommand{
	{Command: "b", Description: "Show current holdings and total balance"},
	{Command: "r", Description: "Remove a holding from your portfolio"},
	{Command: "h", Description: "Show usage instructions"},
	{Command: "start", Description: "Welcome message and reset state"},
}

// New creates a Bot, verifying the token with Telegram.
func New(token string, svc *portfolio.Service, yahoo *finance.YahooClient) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorised as @%s", api.Self.UserName)

	if _, err := api.Request(tgbotapi.NewSetMyCommands(botCommands...)); err != nil {
		log.Printf("set bot commands: %v", err)
	}

	h := newHandler(api, svc, yahoo)
	return &Bot{api: api, handler: h}, nil
}

// Start begins the long-poll update loop. It blocks until ctx is cancelled.
func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			go b.dispatch(ctx, update)
		}
	}
}

// SendMarkdown sends a Markdown-formatted message to a chat.
func (b *Bot) SendMarkdown(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("send markdown to %d: %v", chatID, err)
	}
}

func (b *Bot) dispatch(ctx context.Context, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		b.handler.HandleMessage(ctx, update.Message)
	case update.CallbackQuery != nil:
		b.handler.HandleCallback(ctx, update.CallbackQuery)
	}
}
