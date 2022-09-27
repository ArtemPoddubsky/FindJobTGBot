package app

import (
	log "github.com/sirupsen/logrus"
)

func (a *App) ClearHistory(chatID int64) {
	if err := a.db.ClearHistory(chatID); err != nil {
		log.Errorln("ClearHistory:", err)
	}

	if err := a.SendMessage(chatID, "История просмотренных вакансий была очищена"); err != nil {
		log.Errorln("ClearHistory:", err)
	}
}
