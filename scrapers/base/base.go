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

func Download(raw string, opt string, id int, shref *ScraperHead) string {
	src := baseurlMap[opt] + raw
	log.Println("downloading:", src)

	r, _ := http.Get(src)
	ex := strings.Split(r.Header.Get("Content-Type"), "/")[1]
	dst := fmt.Sprintf("%s/%s/%s_%d.%s",
		storagePath,
		opt,
		shref.Hash,
		id,
		ex)
	log.Println(dst)
	out, _ := os.Create(dst)
	defer r.Body.Close()
	defer out.Close()

	io.Copy(out, r.Body)
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

func (sh *ScraperHead) SHA256() string {
	data := strings.Join([]string{sh.Title, sh.Description, sh.Date}, ":")
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
