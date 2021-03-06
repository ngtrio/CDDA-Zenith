package i18n

import (
	"fmt"
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

func Tran(field string, json *gjson.Result, mo *gotext.Mo) string {
	if raw, has := jsonutil.GetField(field, json); !has || len(raw.String()) == 0 {
		return ""
	} else {
		var res string
		if jsonutil.IsString(raw) {
			res = mo.Get(raw.String())
		} else {
			var str, strSp, ctxt, strPl string
			var has bool
			ctxt, _ = jsonutil.GetString("ctxt", raw, "")
			if str, has = jsonutil.GetString("str", raw, ""); has {
				strPl, _ = jsonutil.GetString("str_pl", raw, "")
			} else if strSp, has = jsonutil.GetString("str_sp", raw, ""); has {
				str = strSp
				strPl = strSp
			}

			if ctxt != "" {
				if strPl != "" {
					res = mo.GetNC(str, strPl, 1, ctxt)
				} else {
					res = mo.GetC(str, ctxt)
				}
			} else {
				if strPl != "" {
					res = mo.GetN(str, strPl, 1)
				} else {
					res = mo.Get(str)
				}
			}
		}
		return res
	}
}

func TranString(raw string, mo *gotext.Mo) string {
	if len(raw) == 0 || mo == nil {
		return raw
	}
	return mo.Get(raw)
}

func TranCustom(raw string, po *gotext.Po, args ...any) string {
	if len(raw) == 0 || po == nil {
		return raw
	}

	td := po.Get(raw)
	if len(args) > 0 {
		return fmt.Sprintf(td, args...)
	}

	return td
}
