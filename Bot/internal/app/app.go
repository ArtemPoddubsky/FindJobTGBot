package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"main/internal/config"
	"main/internal/storage"
	"main/internal/utils"
	"sync"
)

const help = "Начать: название вакансии\nПоиск по последнему запросу: /repeat\nОчистить историю поиска: /clear"

// App holds configuration data, storage and bot instance and map for holding request information.
type App struct {
	config *config.Config
	db     *storage.Postgres
	Bot    *tgbotapi.BotAPI
	Mutex  *sync.RWMutex
	Req    map[int64]string
}

// NewApp returns new instance of App.
func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		db:     storage.NewDB(cfg),
		Bot:    createBot(cfg.Token),
		Mutex:  &sync.RWMutex{},
		Req:    make(map[int64]string, 0),
	}
}

func createBot(token string) *tgbotapi.BotAPI {
	b, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Authorized on account ", b.Self.UserName)

	return b
}

// Run starts updates loop.
func (a *App) Run() {
	log.Infoln("Running")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := a.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			a.Handler(update.Message.Chat.ID, update.Message.Text)
		}
	}
}

// Handler parses requested message and calls corresponding methods.
func (a *App) Handler(id int64, msg string) {
	switch msg {
	case "/start":
		if err := a.sendMessage(id, "Введите название вакансии"); err != nil {
			log.Errorln("/start: ", err)
		}
	case "/repeat":
		a.Repeat(id)
	case "/clear":
		a.ClearHistory(id)
	case "/help":
		if err := a.sendMessage(id, help); err != nil {
			log.Errorln("help:", err)
		}
	case "Не важно", "Нет опыта", "От 1 до 3 лет", "От 3 до 6 лет", "Более 6 лет":
		a.StartSearch(id, msg)
	default:
		a.Mutex.Lock()
		a.Req[id] = msg
		a.Mutex.Unlock()

		if err := a.sendKeyboard(id); err != nil {
			utils.FieldError("Send keyboard:", err, msg)
		}
	}
}
