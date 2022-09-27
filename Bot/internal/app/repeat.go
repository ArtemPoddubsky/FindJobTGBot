package app

import (
	log "github.com/sirupsen/logrus"
)

func (a *App) Repeat(chatID int64) {
	title, exp, err := a.db.GetLast(chatID)

	if err != nil {
		if err.Error() != "no rows in result set" {
			log.Errorln("Repeat:", err)
			return
		}

		if err = a.SendMessage(chatID, "Для того, чтобы пользоваться /repeat нужно ввести ваш первый запрос"); err != nil {
			log.Errorln("Repeat:", err)
		}

		return
	}

	if err = a.Search(chatID, title, exp); err != nil {
		log.Errorln(err)
	}
}
