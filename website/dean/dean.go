package dean

import (
	"JiaoNiBan-scraper/base"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func RequestHRef(url string, page int) ([]base.ScraperHref, int) {

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
				if page == 0 {
					tmp.Href = base.DeanBaseURL + j.Attr("href")[2:]
				} else {
					tmp.Href = base.DeanBaseURL + j.Attr("href")[5:]
				}
				tmp.Title = j.ChildText("h2")
			})
			tmp.Description = i.ChildText("p")
		})
		tmp.Hash = tmp.SHA256()
		tmp.Page = page
		arr = append(arr, tmp)
	})

	var sum int
	c.OnHTML("span[class='p_no']", func(h *colly.HTMLElement) {
		sum, _ = strconv.Atoi(h.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
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

	table_parser := func(h *colly.HTMLElement) string {
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

	img_parser := func(raw string, cnt int) (string, string) {
		tp := strings.Split(raw, ".")[1]
		dst := fmt.Sprintf("%s_%d.%s", shref.Hash, cnt, tp)
		dst = base.ParseWebFile(dst, "dean", shref.Page)
		src := base.DeanBaseURL + raw
		return src, dst
	}
	i_cnt := 1
	e_cnt := 1
	c.OnHTML("div[class='v_news_content']>p,table,p>img", func(h *colly.HTMLElement) {
		if n := h.Name; n == "p" {
			if raw, flag := h.DOM.Attr("src"); flag {
				src, dst := img_parser(raw, i_cnt)
				sc.Text += fmt.Sprintf("([%s])\n", dst)
				base.Download(src, dst)
				i_cnt++
			} else {
				if h.Text != "" {
					sc.Text += h.Text + "\n"
				}
			}
		} else if n == "table" {
			sc.Text += table_parser(h)
		} else if n == "img" {
			raw := h.Attr("src")
			src, dst := img_parser(raw, i_cnt)
			sc.Text += fmt.Sprintf("([%s])\n", dst)
			base.Download(src, dst)
			i_cnt++
		}
	})

	c.OnHTML("div[class='Newslist2']", func(h *colly.HTMLElement) {
		h.ForEach("a[href],table", func(_ int, i *colly.HTMLElement) {

			if n := i.Name; n == "a" {
				href := i.Attr("href")
				tp := strings.Split(href, ".")[1]
				dst := fmt.Sprintf("%s_%d.%s", shref.Hash, e_cnt, tp)
				dst = base.ParseWebFile(dst, "dean", shref.Page)
				src := base.DeanBaseURL + href
				sc.Appendix += fmt.Sprintf("([%s])\n", dst)
				sc.Appendix += i.Text + "\n"
				base.Download(src, dst)
				e_cnt += 1
			} else if n == "table" {
				sc.Appendix += table_parser(i)
			}
		})
	})

	c.Visit(shref.Href)
	return sc
}

func CheckUpdate() {

}

func Request() {

}
