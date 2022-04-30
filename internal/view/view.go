package view

import (
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type View struct {
	Mod     string
	Type    string
	RawJson *gjson.Result
	Mo      *gotext.Mo
	Po      *gotext.Po
}

func (v *View) Render() string {
	tp, _ := jsonutil.GetString("type", v.RawJson, "")
	var obj Type
	switch tp {
	case "MONSTER":
		obj = &VO{
			BaseType: BaseType{
				Mod: v.Mod,
			},
		}
	default:
		log.Warnf("type: %s is not supported to render", tp)
		return ""
	}
	obj.Bind(v.RawJson, v.Mo, v.Po)
	switch v.Type {
	case "cli":
		return obj.CliView(v.Po)
	case "json":
		return obj.JsonView()
	default:
		return ""
	}
}
