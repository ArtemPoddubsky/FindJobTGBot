package utils

import (
	log "github.com/sirupsen/logrus"
)

func FieldError(detail string, err error, request string) {
	var fields = log.Fields{
		"request": request,
	}

	log.WithFields(fields).Errorln(detail, err)
}
