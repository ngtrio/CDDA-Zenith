package main

import (
	"fmt"
	"path"
	"runtime"

	"zenith/internal/core"

	log "github.com/sirupsen/logrus"
)

var game core.Game

func main() {
	for {
		fmt.Print("Zenith> ")
		var input string
		fmt.Scanln(&input)
		res := game.GetById(input, "cli")
		if len(res) == 0 {
			res = game.GetByName(input, "cli")
		}
		for _, out := range res {
			fmt.Println(out)
		}
	}
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

	game = core.Game{
		Mods:    make(map[string]*core.Mod),
		ModPath: "cataclysmdda-0.F/data/mods",
		Lang:    "zh_CN",
	}
	game.Load(map[string]bool{})

	fmt.Println(`
	 ______________________________ 
	/ Hey man! Take Zenith and ME, \
	\ you'll survive!              /
	 ------------------------------ 
    		\   ^__^
    		 \  (oo)\_______
    		    (__)\       )\/\
    		        ||----w |
    		        ||     ||

         - Cataclysm: Dark Days Ahead -
	`)

}
