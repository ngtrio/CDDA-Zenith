package view

import (
	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

type Type interface {
	Bind(json *gjson.Result, mo *gotext.Mo)
	CliView() string
	JsonView() string
}

type BaseType struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Desc   string `json:"description"`
	Symbol string `json:"symbol"`
}

type Monster struct {
	BaseType
	DiffDesc string  `json:"diff_desc"`
	Diff     float64 `json:"difficulty"`
}
