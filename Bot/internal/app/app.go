package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"main/internal/config"
	"main/internal/storage"
	"main/internal/utils"
	"sync"
)

const Help = "Начать: название вакансии\nПоиск по последнему запросу: /repeat\nОчистить историю поиска: /clear"

type App struct {
	config config.Config
	db     *storage.Postgres
	Bot    *tgbotapi.BotAPI
	Mutex  *sync.RWMutex
	Req    map[int64]string
}

func NewApp(cfg config.Config) *App {
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

func (a *App) Handler(id int64, msg string) {
	switch msg {
	case "/start":
		if err := a.SendMessage(id, "Введите название вакансии"); err != nil {
			log.Errorln("Send /start prompt:", err)
		}
		break
	case "/repeat":
		a.Repeat(id)
		break
	case "/clear":
		a.ClearHistory(id)
		break
	case "/help":
		if err := a.SendMessage(id, Help); err != nil {
			log.Errorln("Help:", err)
		}
		break
	case "Не важно", "Нет опыта", "От 1 до 3 лет", "От 3 до 6 лет", "Более 6 лет":
		a.StartSearch(id, msg)
		break
	default:
		a.Mutex.Lock()
		a.Req[id] = msg
		a.Mutex.Unlock()
		if err := a.SendKeyboard(id); err != nil {
			utils.FieldError("Send keyboard:", err, msg)
		}
		break
	}
}
