package jsonutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func Set(json *string, path string, value interface{}) {
	if str, err := sjson.Set(*json, path, value); err != nil {
		log.Error(err)
	} else {
		*json = str
	}
}

func GetString(json *gjson.Result, field string) string {
	f := json.Get(field)
	if f.Exists() {
		return f.String()
	}
	return ""
}
