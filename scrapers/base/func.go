package base

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Download(raw string, opt string, id int, shref *ScraperHead) string {
	src := baseurlMap[opt] + raw

	r, err := http.Get(src)
	if err != nil {
		return "Error Fetching"
	}

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
			s := strings.Split(strings.Split(cd, ";")[2], ".")
			ex = s[len(s)-1]
		}
	}

	dirs := fmt.Sprintf("%s/%s", storagePath, opt)
	static_dirs := fmt.Sprintf("%s/%s", staticPath, opt)
	dst := fmt.Sprintf("%s/%s_%d.%s",
		dirs,
		shref.Hash,
		id,
		ex)
	static_dst := fmt.Sprintf("%s/%s_%d.%s", static_dirs, shref.Hash, id, ex)
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
	return static_dst
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

func SH_SHA256(sh ScraperHead) string {
	data := strings.Join([]string{sh.Title, sh.Desc, strconv.Itoa(sh.Year),
		strconv.Itoa(sh.Month), strconv.Itoa(sh.Day)}, ":")
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
