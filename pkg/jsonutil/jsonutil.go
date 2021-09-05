package jsonutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"
)

func Set(json *string, path string, value interface{}) {
	if str, err := sjson.Set(*json, path, value); err != nil {
		log.Error(err)
	} else {
		*json = str
	}
}
