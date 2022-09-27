package app

import (
	log "github.com/sirupsen/logrus"
)

// ClearHistory calls database to clear user request history
// and sends message in case of success.
func (a *App) ClearHistory(chatID int64) {
	if err := a.db.ClearHistory(chatID); err != nil {
		log.Errorln("ClearHistory:", err)
	}

	if err := a.sendMessage(chatID, "История просмотренных вакансий была очищена"); err != nil {
		log.Errorln("ClearHistory:", err)
	}
}
