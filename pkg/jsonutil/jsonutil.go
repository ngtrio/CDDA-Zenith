package jsonutil

import (
	"encoding/json"
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

func GetArray(field string, json *gjson.Result, df []gjson.Result) ([]gjson.Result, bool) {
	if res, has := GetField(field, json); has {
		return res.Array(), has
	} else {
		return df, false
	}
}

func GetInt(field string, json *gjson.Result, df int64) (int64, bool) {
	if res, has := GetField(field, json); has {
		return res.Int(), has
	} else {
		return df, false
	}
}

func GetFloat(field string, json *gjson.Result, df float64) (float64, bool) {
	if res, has := GetField(field, json); has {
		return res.Float(), has
	} else {
		return df, false
	}
}

func GetString(field string, json *gjson.Result, defaultValue string) (string, bool) {
	if res, has := GetField(field, json); has {
		return res.String(), has
	} else {
		return defaultValue, false
	}
}

func GetField(field string, json *gjson.Result) (*gjson.Result, bool) {
	res := json.Get(field)
	if res.Exists() {
		return &res, true
	} else {
		return nil, false
	}
}

func IsString(json *gjson.Result) bool {
	return json.Type == gjson.String
}

func ToJson(s any) string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
