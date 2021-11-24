package base

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
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
	Page int
	Href string
}

type ScraperContent struct {
	ScraperHead
	Text     string
	Appendix string
}

func Download(src string, dst string) {

}

func ParseWebFile(dst string, cat string, page int) string {
	return fmt.Sprintf("/downloads/website/%s/%d/%s", cat, page, dst)
}

func (sh *ScraperHead) SHA256() string {
	data := strings.Join([]string{sh.Title, sh.Description, sh.Date}, ":")
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}
