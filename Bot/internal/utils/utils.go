package utils

import (
	log "github.com/sirupsen/logrus"
)

const Help = "Начать: название вакансии\nПоиск по последнему запросу: /repeat\nОчистить историю поиска: /clear"
const Repeat = "Используйте /repeat чтобы посмотреть следующие 7 вакансий"

func Error(err error, request ...string) {
	if request == nil {
		log.Errorln(err)
	} else {
		log.WithFields(log.Fields{
			"request": request[0],
		}).Errorln(err)
	}
}
