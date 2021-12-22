package main

import (
	"JiaoNiBan-data/scrapers/website/dean"
	"flag"
	"log"
)

func main() {
	pages := flag.Int("s", -1, "Setup the Scraper")
	upd := flag.Bool("c", false, "Checking Data Updates")
	flag.Parse()
	if *pages != -1 {
		log.Printf("Fetching %d pages for database", *pages)
		dean.Setup(*pages)
		return
	}
	if *upd {
		log.Println("Checking Update for Dean Office")
		dean.CheckUpdate()
		return
	}
	// var sh base.ScraperHref
	// sh.Href = "https://www.sjtu.edu.cn/tg/20211220/165433.html"
	// fmt.Println(dean.RequestOne(sh, 0))
}
