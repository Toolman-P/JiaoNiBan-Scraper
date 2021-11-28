package dean

import (
	"JiaoNiBan-data/databases"
	"JiaoNiBan-data/scrapers/base"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func parseIndex(index int) string {
	return fmt.Sprintf("%s/%d.htm", base.DeanPrefix, index)
}

func requestHRef(url string, mode int) ([]base.ScraperHref, int) {

	var arr []base.ScraperHref
	c := colly.NewCollector(
		colly.UserAgent(base.UserAgent),
		colly.MaxDepth(5),
	)

	c.OnHTML("li[class='clearfix']", func(h *colly.HTMLElement) {

		tmp := base.ScraperHref{}
		tmp.Author = "dean"
		h.ForEach("div[class='sj']", func(_ int, i *colly.HTMLElement) {
			day := i.ChildText("h2")
			year_month := i.ChildText("p")
			tmp.Date = strings.Join([]string{year_month, day}, ".")
		})

		h.ForEach("div[class='wz']", func(_ int, i *colly.HTMLElement) {
			i.ForEach("a[href]", func(_ int, j *colly.HTMLElement) {
				if mode == 0 {
					tmp.Href = base.DeanBaseURL + j.Attr("href")[2:]
				} else {
					tmp.Href = base.DeanBaseURL + j.Attr("href")[5:]
				}
				tmp.Title = j.ChildText("h2")
			})
			tmp.Description = i.ChildText("p")
		})
		tmp.Hash = base.SHA256(tmp.ScraperHead)
		arr = append(arr, tmp)
	})

	var sum int
	c.OnHTML("span[class='p_no']", func(h *colly.HTMLElement) {
		sum, _ = strconv.Atoi(h.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)
	})

	c.Visit(url)

	return arr, sum
}

func requestOne(shref base.ScraperHref) base.ScraperContent {

	sc := base.ScraperContent{}
	sc.ScraperHead = shref.ScraperHead
	sc.Text = ""
	sc.Appendix = ""

	c := colly.NewCollector(
		colly.MaxDepth(5),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", base.UserAgent)
		r.Headers.Set("Referer", base.DeanBaseURL)
	})

	i_cnt := 1
	c.OnHTML("div[class='v_news_content']>p,table,p>img", func(h *colly.HTMLElement) {

		if n := h.Name; n == "p" {
			if raw, flag := h.DOM.Attr("src"); flag {
				dst := base.Download(raw, "dean", i_cnt, &shref.ScraperHead)
				sc.Text += fmt.Sprintf("([%s])\n", dst)
				i_cnt++
			} else {
				if h.Text != "" {
					sc.Text += h.Text + "\n"
				}
			}
		} else if n == "table" {
			sc.Text += base.ParseTable(h)
		} else if n == "img" {
			raw := h.Attr("src")
			dst := base.Download(raw, "dean", i_cnt, &shref.ScraperHead)
			sc.Text += fmt.Sprintf("([%s])\n", dst)
			i_cnt++
		}
	})

	e_cnt := 1
	c.OnHTML("div[class='Newslist2']", func(h *colly.HTMLElement) {

		h.ForEach("a[href],table", func(_ int, i *colly.HTMLElement) {

			if n := i.Name; n == "a" {
				raw := i.Attr("href")
				dst := base.Download(raw, "dean", e_cnt, &shref.ScraperHead)
				sc.Appendix += i.Text + "\n"
				sc.Appendix += fmt.Sprintf("([%s])\n", dst)

				e_cnt += 1
			} else if n == "table" {
				sc.Appendix += base.ParseTable(i)
			}
		})
	})
	c.Visit(shref.Href)
	return sc
}

func requestContents(shrefs []base.ScraperHref) {
	var w sync.WaitGroup
	for _, h := range shrefs {
		if f, _ := databases.CheckHrefExists("dean", h.Hash); !f {
			w.Add(1)
			go func(i base.ScraperHref) {
				databases.AddHref("dean", i.Hash)
				c := requestOne(i)
				databases.AddPage(c)
				w.Done()
			}(h)
		}
	}
	w.Wait()
}

func validateVersion() bool {
	r, _ := http.Get(base.DeanBaseURL)
	t, _ := ioutil.ReadAll(r.Body)
	h := sha256.Sum256(t)
	if sha := hex.EncodeToString(h[:]); sha != databases.GetVersion("dean") {
		databases.SetVersion("dean", sha)
		return false
	}
	return true
}

func CheckUpdate() {

	databases.Init()
	defer databases.Close()

	if !validateVersion() {
		log.Println("Fetching updates...")
		hrefs, _ := requestHRef(base.DeanFirstPage, 0)
		requestContents(hrefs)
		return
	}
	log.Println("Current version is the latest version.")
}

func Setup(pages int) {
	databases.Init()
	defer databases.Close()
	validateVersion()
	fhref, sum := requestHRef(base.DeanFirstPage, 0)
	requestContents(fhref)
	for i := sum - 1; i >= sum-pages; i-- {
		hrefs, _ := requestHRef(parseIndex(i), 1)
		requestContents(hrefs)
	}
}
