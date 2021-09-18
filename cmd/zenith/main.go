package main

import (
	"fmt"
	"path"
	"runtime"

	"zenith/internal/core"
	"zenith/internal/data"

	log "github.com/sirupsen/logrus"
)

func main() {
	// mo := loader.LoadLang("zh_CN")
	// str := mo.Get("caustic soldier zombie")
	// log.Info(str)
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
	log.SetReportCaller(true)
	// log.SetLevel(log.DebugLevel)

	game := core.Game{
		Mods:    make(map[string]*data.Mod),
		ModPath: "cataclysmdda-0.F/data/mods",
		Lang:    "zh_CN",
	}
	game.Load(map[string]bool{"dda": true})

	res := game.GetById("mon_zombie_kevlar_2", "cli")
	res1 := game.GetById("mon_zombie_hulk", "cli")
	log.Info(res, res1)

}
