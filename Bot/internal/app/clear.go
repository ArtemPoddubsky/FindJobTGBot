package app

import (
	"fmt"
	"main/internal/utils"
)

func (a *App) ClearHistory(chatID int64) {
	if err := a.db.ClearHistory(chatID); err != nil {
		utils.Error(fmt.Errorf("ClearHistory: %v ", err))
	} else if err = a.SendMessage(chatID, "История просмотренных вакансий была очищена"); err != nil {
		utils.Error(fmt.Errorf("ClearHistory: %v ", err))
	}
}
