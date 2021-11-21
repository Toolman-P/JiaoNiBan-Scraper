package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type HeadInfo struct {
	title       string
	description string
	date        string
	href        string
}

type Contents struct {
}

type NotifMap map[string]HeadInfo
type ContentMap map[string]string

func get_sha256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func get_hrefmap(user_agent string, url string) NotifMap {

	m := make(NotifMap)

	c := colly.NewCollector(
		colly.UserAgent(user_agent),
		colly.MaxDepth(5),
	)

	c.OnHTML("li[class='clearfix']", func(e1 *colly.HTMLElement) {

		var date string
		e1.ForEach("div[class='sj']", func(_ int, e2 *colly.HTMLElement) {
			day := e2.ChildText("h2")
			year_month := e2.ChildText("p")
			date = strings.Join([]string{year_month, day}, ".")
		})

		var href string
		var title string
		var description string
		e1.ForEach("div[class='wz']", func(_ int, e2 *colly.HTMLElement) {
			e2.ForEach("a[href]", func(_ int, e3 *colly.HTMLElement) {
				href = e3.Attr("href")
				title = e3.ChildText("h2")
			})
			description = e2.ChildText("p")
		})
		h := get_sha256(strings.Join([]string{title, description, date, href}, ":"))
		m[h] = HeadInfo{title, description, date, href}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)

	return m
}

func get_content(user_agent string, referer string) {

	c := colly.NewCollector(
		colly.MaxDepth(5),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", user_agent)
		r.Headers.Set("Referer", referer)
	})

	table_parser := func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, i *colly.HTMLElement) {
			fmt.Print("[")
			i.ForEach("td", func(_ int, j *colly.HTMLElement) {
				if j.Text == "" {
					fmt.Print("/ ")
				} else {
					fmt.Print(j.Text + " ")
				}
			})
			fmt.Println("]")
		})
	}

	c.OnHTML("div[class='v_news_content']>p,table", func(e *colly.HTMLElement) {
		if n := e.Name; n == "p" {
			if src, flag := e.DOM.Attr("img"); flag {
				fmt.Println(src)
			}
			fmt.Println(e.Text)
		} else if n == "table" {
			table_parser(e)
		}
	})

	c.OnHTML("div[class='Newslist2']", func(e *colly.HTMLElement) {
		e.ForEach("a[href],table", func(_ int, h *colly.HTMLElement) {
			if n := h.Name; n == "a" {
				href := h.Attr("href")
				fmt.Println(href)
				fmt.Println(h.Text)
			} else if n == "table" {
				table_parser(h)
			}
		})
	})

	// c.OnHTML("div[class='Newslist2']", func(e1 *colly.HTMLElement) {
	// 	e1.DOM.Each()
	// })
	c.Visit("https://jwc.sjtu.edu.cn/info/1222/11297.htm")
	// c.Visit("https://jwc.sjtu.edu.cn/info/1222/11294.htm")
}

func main() {
	// m := HrefMap{}
	user_agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36"
	notif_url := "https://jwc.sjtu.edu.cn/index/mxxsdtz.htm"
	// get_hrefmap(user_agent, notif_url)
	get_content(user_agent, notif_url)
}
