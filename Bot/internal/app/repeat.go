package app

import (
	log "github.com/sirupsen/logrus"
)

// Repeat executes last request for specific user.
func (a *App) Repeat(chatID int64) {
	title, exp, err := a.db.GetLast(chatID)

	if err != nil {
		if err.Error() != "no rows in result set" {
			log.Errorln("repeat:", err)
			return
		}

		if err = a.sendMessage(chatID, "Для того, чтобы пользоваться /repeat нужно ввести ваш первый запрос"); err != nil {
			log.Errorln("repeat:", err)
		}

		return
	}

	if err = a.search(chatID, title, exp); err != nil {
		log.Errorln(err)
	}
}
