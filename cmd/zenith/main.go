package main

import (
	"fmt"
	"path"
	"runtime"

	"zenith/internal/core"
	"zenith/internal/data"

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
		Mods:    make(map[string]*data.Mod),
		ModPath: "cataclysmdda-0.F/data/mods",
		Lang:    "zh_CN",
	}
	game.Load(map[string]bool{})

	fmt.Println(`
	 __________ _   _ ___ _____ _   _ 
	|__  / ____| \ | |_ _|_   _| | | |
	  / /|  _| |  \| || |  | | | |_| |
	 / /_| |___| |\  || |  | | |  _  |
	/____|_____|_| \_|___| |_| |_| |_|

	  - Cataclysm: Dark Days Ahead -													
	`)
}
