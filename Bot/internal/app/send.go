package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"main/internal/app/parser"
)

func (a *App) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := a.Bot.Send(msg)
	return err
}

func (a *App) SendKeyboard(chatID int64) error {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Нет опыта"),
			tgbotapi.NewKeyboardButton("От 1 до 3 лет")),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("От 3 до 6 лет"),
			tgbotapi.NewKeyboardButton("Более 6 лет")),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Не важно")),
	)
	msg := tgbotapi.NewMessage(chatID, "Выберите ваш опыт работы")
	msg.ReplyMarkup = keyboard
	_, err := a.Bot.Send(msg)
	return err
}

func (a *App) SendVacancy(chatID int64, vacancy parser.Vacancy, salary string) error {
	msg := tgbotapi.NewMessage(chatID, vacancy.Employer.Name+" . "+vacancy.Name+"\n"+salary+"\n"+vacancy.Url+"\n")
	_, err := a.Bot.Send(msg)
	return err
}
