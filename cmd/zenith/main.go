package main

import (
	"fmt"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"

	"zenith/internal/core"
	"zenith/internal/data"
)

func main() {

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
		ModPath: "data/mods",
	}
	game.LoadMod(map[string]bool{"dda": true})
}
