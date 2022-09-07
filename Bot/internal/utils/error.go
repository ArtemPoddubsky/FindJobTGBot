package utils

import (
	log "github.com/sirupsen/logrus"
)

func Error(err error, request ...string) {
	if request == nil {
		log.Errorln(err)
	} else {
		log.WithFields(log.Fields{
			"request": request[0],
		}).Errorln(err)
	}
}
