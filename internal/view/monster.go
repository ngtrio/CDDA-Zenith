package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"zenith/internal/color"
	"zenith/internal/i18n"
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

func (m *Monster) Bind(raw *gjson.Result, mo *gotext.Mo) {
	m.ID, _ = jsonutil.GetString("id", raw, "")
	m.Type, _ = jsonutil.GetString("type", raw, "")
	m.Name = i18n.Tran("name", raw, mo)
	m.Desc = i18n.Tran("description", raw, mo)

	cr, _ := jsonutil.GetString("color", raw, "")
	symbol, _ := jsonutil.GetString("symbol", raw, "")
	var colorLoader color.Color
	colorLoader.Load(cr)
	m.Symbol = colorLoader.Colorized(symbol)

	diffColor, _ := jsonutil.GetString("diff_color", raw, "")
	temp := i18n.Tran("diff_desc", raw, mo)
	l := strings.Index(temp, ">")
	r := strings.LastIndex(temp, "<")
	diffDesc := temp[l+1 : r]
	colorLoader.Load(diffColor)
	m.DiffDesc = colorLoader.Colorized(diffDesc)

	value, _ := jsonutil.GetFloat("difficulty", raw, 0)
	m.Diff, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
}

func (m *Monster) CliView() string {
	template := `
%s %s
%s(%.3f)
%s`
	return fmt.Sprintf(template, m.Symbol, m.Name, m.DiffDesc, m.Diff, m.Desc)
}

func (m *Monster) JsonView() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}
