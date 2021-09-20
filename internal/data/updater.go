package data

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"
	"zenith/internal/config"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

func UpdateNow(useProxy bool) bool {
	url, err := getLatestUrl()
	if err != nil {
		log.Error(err)
		return false
	}
	if useProxy {
		url = config.GHProxy + url
	}

	if resp, err := http.Get(url); err != nil {
		log.Error(err)
		return false
	} else {
		if f, err := os.OpenFile(config.DownloadPath, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			log.Error(err)
			return false
		} else {
			bar := progressbar.DefaultBytes(resp.ContentLength, "downloading")
			io.Copy(io.MultiWriter(f, bar), resp.Body)
			return deCompress()
		}
	}
}

func getLatestUrl() (string, error) {
	url := config.ReleasePage
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	latest := doc.Find("h2").First()

	list := latest.SiblingsFiltered("ul").First()
	var link string
	list.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		name := s.Text()
		if strings.Contains(name, "linux-tiles") {
			link, _ = s.Attr("href")
			return false
		}
		return true
	})

	return link, nil
}

func deCompress() bool {
	srcFile, err := os.Open(config.DownloadPath)
	if err != nil {
		log.Error(err)
		return false
	}
	defer srcFile.Close()

	gz, err := gzip.NewReader(srcFile)
	if err != nil {
		log.Error(err)
		return false
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return true
		}
		if err != nil {
			log.Error(err)
			return false
		}
		idx := strings.Index(hdr.Name, "/")
		filename := config.BaseDir + hdr.Name[idx:]
		f, err := creatFileOrDir(filename)
		if err != nil {
			log.Error(err)
			return false
		}
		if f != nil {
			io.Copy(f, tr)
		}
	}
}

func creatFileOrDir(path string) (*os.File, error) {
	_, err := os.Stat(path)

	if os.IsExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		idx := strings.LastIndex(path, "/")
		err = os.MkdirAll(path[0:idx], 0755)
		if err != nil {
			return nil, err
		}
		if idx != len(path)-1 {
			return os.Create(path)
		}
	}
	return nil, err
}
