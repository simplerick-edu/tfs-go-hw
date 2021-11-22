package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	chats map[int64]*tgbotapi.Chat
	bot   *tgbotapi.BotAPI
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	return &Bot{
		chats: make(map[int64]*tgbotapi.Chat),
		bot:   bot,
	}, err
}

func (t *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates, err := t.bot.GetUpdatesChan(updateConfig)
	go func() {
		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}
			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help":
				msg.Text = "I understand /start and /stop"
			case "start":
				t.addChat(update.Message.Chat)
				msg.Text = "You have successfully subscribed to notifications"
			case "stop":
				t.removeChat(update.Message.Chat.ID)
				msg.Text = "You have successfully unsubscribed from notifications"
			default:
				msg.Text = "See /help for list of commands"
			}

			if _, err := t.bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		}
	}()
	return err
}

func (t *Bot) Stop() {
	t.bot.StopReceivingUpdates()
}

func (t *Bot) addChat(chat *tgbotapi.Chat) {
	t.chats[chat.ID] = chat
}

func (t *Bot) removeChat(chatID int64) {
	delete(t.chats, chatID)
}

func (t *Bot) Notify(text string) error {
	msg := tgbotapi.NewMessage(0, text)
	for chatID := range t.chats {
		msg.ChatID = chatID
		_, err := t.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	fmt.Println("notification sent")
	return nil
}
