package loader

import (
	"zenith/internal/config"

	"github.com/leonelquinteros/gotext"
)

func LoadLang(lang string) *gotext.Mo {
	path := config.BaseDir + "/lang/mo/" + lang + "/LC_MESSAGES/cataclysm-dda.mo"
	mo := gotext.NewMo()
	mo.ParseFile(path)
	return mo
}
