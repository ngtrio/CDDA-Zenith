package loader

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
	"zenith/pkg/fileutil"

	"github.com/tidwall/gjson"
)

// load data from paths
func LoadJsonFromPaths(paths ...string) []*gjson.Result {
	var jsons []*gjson.Result
	for _, path := range paths {
		if files, dirs, err := fileutil.Ls(path); err != nil {
			log.Fatalf("read dir %s fail, err: %v", path, err)
		} else {
			sortPaths(files)
			sortPaths(dirs)
			// load normal json file
			for _, file := range files {
				res := LoadJsonFromFile(file)
				jsons = append(jsons, res...)

			}

			for _, dir := range dirs {
				res := LoadJsonFromPaths(dir)
				jsons = append(jsons, res...)
			}
		}
	}
	return jsons
}

func LoadJsonFromFile(file string) []*gjson.Result {
	if !strings.HasSuffix(file, ".json") {
		return nil
	}

	if bytes, err := fileutil.ReadFile(file); err != nil {
		log.Fatalf("read file %s fail, err: %v\n", file, err)
		return nil
	} else {
		r := loadJson(bytes)
		log.Debugf("[LOADED]: %s\n", file)
		return r
	}
}

func loadJson(bytes []byte) []*gjson.Result {
	json := gjson.ParseBytes(bytes)
	res := make([]*gjson.Result, 0)
	if json.IsArray() {
		json.ForEach(func(k, v gjson.Result) bool {
			res = append(res, &v)
			return true
		})
	} else if json.IsObject() {
		res = append(res, &json)
	}
	return res
}

func sortPaths(paths []string) {
	sort.SliceStable(paths, func(i, j int) bool {
		aSplits := strings.Split(paths[i], "/")
		bSplits := strings.Split(paths[j], "/")
		return aSplits[len(aSplits)-1] <= bSplits[len(bSplits)-1]
	})
}

func NeedInherit(json *gjson.Result) bool {
	return json.Get("copy-from").Exists()
}

func IsAbstract(json *gjson.Result) bool {
	return json.Get("abstract").Exists()
}
