package base

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper struct {
	user_agent string
	referer    string
	base_url   string
	target     string
}

type ScraperHead struct {
	Title       string
	Author      string
	Date        string
	Description string
	Hash        string
}

type ScraperHref struct {
	ScraperHead
	Href string
}

type ScraperContent struct {
	ScraperHead
	Text     string
	Appendix string
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Download(raw string, opt string, id int, shref *ScraperHead) string {
	src := baseurlMap[opt] + raw

	r, _ := http.Get(src)
	log.Println("Downloading:", src)
	var ex string
	{
		cd := r.Header.Get("Content-Disposition")
		if len(cd) == 0 {
			ct := r.Header.Get("Content-Type")
			if len(ct) == 0 {
				return "Failed Fetching"
			}
			ex = strings.Split(ct, "/")[1]

		} else {
			s1 := strings.Split(cd, ";")
			ex = strings.Split(s1[2], ".")[1]
		}
	}

	dirs := fmt.Sprintf("%s/%s", storagePath, opt)
	dst := fmt.Sprintf("%s/%s_%d.%s",
		dirs,
		shref.Hash,
		id,
		ex)
	log.Println("Saved to:", dst)

	defer r.Body.Close()
	if !isExist(dirs) {
		os.MkdirAll(dirs, os.ModePerm)
	}
	if !isExist(dst) {
		out, _ := os.Create(dst)
		defer out.Close()
		io.Copy(out, r.Body)
	}
	return dst
}

func ParseTable(h *colly.HTMLElement) string {
	var rt string
	h.ForEach("tr", func(_ int, i *colly.HTMLElement) {
		rt += "[["
		i.ForEach("td", func(_ int, j *colly.HTMLElement) {
			if j.Text == "" {
				rt += "/ "
			} else {
				rt += j.Text + " "
			}
		})
		rt += "]]\n"
	})
	return rt
}

func SHA256(sh ScraperHead) string {
	data := strings.Join([]string{sh.Title, sh.Description, sh.Date}, ":")
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
