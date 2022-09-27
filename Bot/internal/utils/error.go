package utils

import (
	log "github.com/sirupsen/logrus"
)

// FieldError performs logging with fields.
func FieldError(detail string, err error, request string) {
	var fields = log.Fields{
		"request": request,
	}

	log.WithFields(fields).Errorln(detail, err)
}
