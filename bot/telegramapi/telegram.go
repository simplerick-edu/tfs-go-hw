package telegramapi

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

var JobQueueIsFull = errors.New("notifier's job queue is full")

const DefaultJobQueueSize = 10

type TgBot struct {
	chatsMu      sync.Mutex
	chats        map[int64]*tgbotapi.Chat
	bot          *tgbotapi.BotAPI
	jobQueue     chan string
	jobQueueSize int
}

func New(token string, jobQueueSize int) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	return &TgBot{
		chats:        make(map[int64]*tgbotapi.Chat),
		bot:          bot,
		jobQueueSize: jobQueueSize,
	}, err
}

func NewWithCreds(token string) (*TgBot, error) {
	return New(token, DefaultJobQueueSize)
}

func (t *TgBot) Start() error {
	// updates handling
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
				log.Println(err)
			}
		}
	}()
	// job queue: listen and send messages
	t.jobQueue = make(chan string, t.jobQueueSize)
	go func() {
		for message := range t.jobQueue {
			err := t.sendMessage(message)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	return err
}

func (t *TgBot) Stop() {
	close(t.jobQueue)
	t.bot.StopReceivingUpdates()
}

func (t *TgBot) addChat(chat *tgbotapi.Chat) {
	t.chatsMu.Lock()
	defer t.chatsMu.Unlock()
	t.chats[chat.ID] = chat
}

func (t *TgBot) removeChat(chatID int64) {
	t.chatsMu.Lock()
	defer t.chatsMu.Unlock()
	delete(t.chats, chatID)
}

func (t *TgBot) sendMessage(text string) error {
	t.chatsMu.Lock()
	defer t.chatsMu.Unlock()
	msg := tgbotapi.NewMessage(0, text)
	msg.ParseMode = "markdown"
	for chatID := range t.chats {
		msg.ChatID = chatID
		_, err := t.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	if len(t.chats) > 0 {
		log.Println("notification sent")
	}
	return nil
}

func (t *TgBot) Notify(text string) error {
	select {
	case t.jobQueue <- text:
		return nil
	default:
		return JobQueueIsFull
	}
}
