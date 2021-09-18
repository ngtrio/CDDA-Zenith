package loader

import (
	"zenith/internal/config"

	"github.com/leonelquinteros/gotext"
)

func LoadMo(lang string) *gotext.Mo {
	path := config.BaseDir + "/lang/mo/" + lang + "/LC_MESSAGES/cataclysm-dda.mo"
	mo := gotext.NewMo()
	mo.ParseFile(path)
	return mo
}

func LoadPo(lang string) *gotext.Po {
	path := "./lang/po/" + lang + ".po"
	po := gotext.NewPo()
	po.ParseFile(path)
	return po
}
