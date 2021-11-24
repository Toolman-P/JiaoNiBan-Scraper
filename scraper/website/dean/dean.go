package dean

import (
	"JiaoNiBan-data/scraper/base"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func parseIndex(i int) {

}

func RequestHRef(url string, mode int) ([]base.ScraperHref, int) {

	var arr []base.ScraperHref
	c := colly.NewCollector(
		// colly.Async(true),
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
		tmp.Hash = tmp.SHA256()
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

func RequestContent(shref base.ScraperHref) base.ScraperContent {

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

	c.OnHTML("div[class='v_news_content']>p,table,p>img", func(h *colly.HTMLElement) {
		i_cnt := 1
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

	c.OnHTML("div[class='Newslist2']", func(h *colly.HTMLElement) {
		e_cnt := 1
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

func Request() {
	hrefs, _ := RequestHRef(base.DeanFirstPage, 0)

	for _, h := range hrefs {
		RequestContent(h)
	}
}
