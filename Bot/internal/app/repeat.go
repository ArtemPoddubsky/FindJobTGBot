package app

import (
	"fmt"
	"main/internal/utils"
)

func (a *App) Repeat(chatID int64) {
	title, exp, err := a.db.GetLast(chatID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			if err = a.SendMessage(chatID, "Для того, чтобы пользоваться /repeat нужно ввести ваш первый запрос"); err != nil {
				utils.Error(fmt.Errorf("Repeat: %v ", err))
			}
		} else {
			utils.Error(fmt.Errorf("Repeat: %v ", err))
		}
		return
	}
	if err = a.Search(chatID, title, exp); err != nil {
		utils.Error(err)
	}
}
