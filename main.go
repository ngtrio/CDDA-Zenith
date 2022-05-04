package main

import (
	"bufio"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leonelquinteros/gotext"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"zenith/internal/view"

	"zenith/internal/config"
	"zenith/internal/core"
	"zenith/internal/data"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var game atomic.Value

func main() {

	options, lang := readOptions()

	if options["--help"] {
		printHelp()
		return
	}

	if !options["--disable-banner"] {
		printBanner()
	}

	configLog(options["--debug-mode"])
	download(options["--use-proxy"], options["--update-now"])
	loadData(getVersion())
	bgTask(options["--use-proxy"])

	if options["--web-mode"] {
		web()
	} else {
		cli(lang)
	}
}

func readOptions() (map[string]bool, string) {
	args := os.Args
	fmt.Printf("options: %v\n", args)

	res := map[string]bool{
		"--help":           false,
		"--use-proxy":      false,
		"--debug-mode":     false,
		"--update-now":     false,
		"--disable-banner": false,
		"--web-mode":       false,
	}
	lang := "zh_CN"
	for _, arg := range args[1:] {
		if _, has := res[arg]; has {
			res[arg] = true
		} else {
			if strings.HasPrefix(arg, "--lang") {
				parts := strings.Split(arg, ":")
				if len(parts) != 2 {
					fmt.Println("[WARN] Language option is invalid, fallback to use zh_CN")
				} else {
					lang = parts[1]
				}
			}
		}
	}

	return res, lang
}

func printBanner() {
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

func printHelp() {
	fmt.Println("Usage: zenith [--help] [--use-proxy] [--debug-mode] [--update-now] [--disable-banner]")
}

func configLog(debug bool) {
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
	// log.SetReportCaller(true)
	if debug {
		log.SetLevel(log.DebugLevel)
		fmt.Println("Home mode is enabled")
	}
}

func download(useProxy, useLatest bool) bool {
	if useLatest {
		if useProxy {
			fmt.Printf("Use proxy to download game data, thanks to: %s\n", config.GHProxy)
		}
		return data.UpdateNow(useProxy)
	}

	_, err := os.Stat(config.BaseDir)
	if err != nil {
		_, err := os.Stat(config.DownloadPath)
		if err != nil {
			fmt.Println("Game data not found, download")
			return data.UpdateNow(useProxy)
		} else {
			fmt.Println("Game data found compressed, decompress")
			return data.DeCompress()
		}
	}

	return true
}

func getVersion() string {
	f, err := os.Open(config.BaseDir + "/VERSION.txt")
	if err != nil && os.IsNotExist(err) {
		fmt.Println("Game data is not found, try to use '--update-now' option.")
		os.Exit(1)
	}

	bytes, _ := ioutil.ReadAll(f)
	return string(bytes)
}

func loadData(version string) {
	g := &core.Game{
		Commit:    strings.Split(strings.Split(version, "\n")[2], ": ")[1][:8],
		Mods:      make(map[string]*core.Mod),
		ModPath:   config.BaseDir + "/data/mods",
		LangPacks: make(map[string]core.LangPack),
	}
	g.Load(map[string]bool{})
	g.UpdateAt = time.Now().Format("2006-01-02 15:04:05.000 -07")

	fmt.Printf("Game version:\n%s\n", version)

	game.Store(g)
}

func cli(lang string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Zenith> ")

		var input string
		if scanner.Scan() {
			input = scanner.Text()
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Bye!")
			os.Exit(0)
		}

		res := getGame().Indexer.IdIndex("MONSTER", input, lang)
		if len(res) == 0 {
			res = getGame().Indexer.NameIndex("MONSTER", input, lang)
		}
		for _, out := range res {
			fmt.Println(out.CliView(getGame().LangPacks[lang].Po))
		}
	}
}

func web() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Static("web"))
	e.Use(middleware.RateLimiterWithConfig(config.NewRateLimiterConfig()))

	e.Renderer = view.NewTemplate()

	e.GET("/", Home)
	e.GET("/info/:mod/:type/:idName", Detail)
	e.GET("/list", List)
	e.GET("/search", Search)

	e.Logger.Fatal(e.Start(":1323"))
}

func Home(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/list?type=MONSTER&num=20&page=1")
}

func Detail(c echo.Context) error {
	modId := c.Param("mod")
	tp := c.Param("type")
	idName := c.Param("idName")
	lang := c.QueryParam("lang")
	if _, has := getGame().LangPacks[lang]; !has {
		lang = "zh_CN"
	}

	res := getGame().Indexer.ModIdIndex(modId, tp, idName, lang)
	if res == nil {
		res = getGame().Indexer.ModNameIndex(modId, tp, idName, lang)
	}

	return c.Render(http.StatusOK, "detail", wrapParam(c, res, lang))
}

func List(c echo.Context) error {
	numParam := c.QueryParam("num")
	pageParam := c.QueryParam("page")
	typeParam := c.QueryParam("type")
	lang := c.QueryParam("lang")
	if _, has := getGame().LangPacks[lang]; !has {
		lang = "zh_CN"
	}

	num, _ := strconv.ParseInt(numParam, 10, 32)
	page, _ := strconv.ParseInt(pageParam, 10, 32)

	if num <= 0 {
		num = 10
	}

	if page <= 0 {
		page = 1
	}

	res, totalPage := getGame().Indexer.RangeIndex(typeParam, int(num), int(page), lang)
	return c.Render(http.StatusOK, "list", wrapParam(c, genListParam(res, int(page), totalPage, numParam, typeParam), lang))
}

func Search(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	tp := c.QueryParam("type")
	lang := c.QueryParam("lang")
	if _, has := getGame().LangPacks[lang]; !has {
		lang = "zh_CN"
	}

	if len(keyword) <= 0 {
		return c.Render(http.StatusOK, "search", wrapParam(c, nil, lang))
	}

	res := getGame().Indexer.FuzzyNameIndex(tp, keyword, lang)
	return c.Render(http.StatusOK, "search", wrapParam(c, tableParam{Items: res}, lang))
}

type tableParam struct {
	Items []*core.VO
}

type listParam struct {
	tableParam
	Type      string
	Num       string
	CurPage   int
	TotalPage int
	NextPage  int
	PrevPage  int
}

func genListParam(items []*core.VO, curPage, totalPage int, num, type_ string) listParam {
	return listParam{
		tableParam: tableParam{
			Items: items,
		},
		Type:      type_,
		Num:       num,
		CurPage:   curPage,
		TotalPage: totalPage,
		NextPage:  curPage + 1,
		PrevPage:  curPage - 1,
	}
}

type templateParam struct {
	Path     string
	Commit   string
	UpdateAt string
	Po       *gotext.Po
	Data     any
}

func wrapParam(c echo.Context, data any, lang string) templateParam {
	return templateParam{
		Path:     c.Path(),
		Commit:   getGame().Commit,
		UpdateAt: getGame().UpdateAt,
		Po:       getGame().LangPacks[lang].Po,
		Data:     data,
	}
}

func getGame() *core.Game {
	return game.Load().(*core.Game)
}

func bgTask(useProxy bool) {
	c := cron.New(cron.WithSeconds())
	id, err := c.AddFunc("0 0 0 * * *", func() {
		download(useProxy, true)
		loadData(getVersion())
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("job id: %v\n", id)
	c.Start()
}
