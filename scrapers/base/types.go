package base

type ScraperHead struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"-"`
	Month  int    `json:"-"`
	Day    int    `json:"-"`
	Date   string `json:"date"`
	Desc   string `json:"desc"`
	Hash   string `json:"-"`
}

type ScraperHref struct {
	ScraperHead
	Href string `json:"href"`
}

type ScraperContent struct {
	ScraperHead
	Body     string
	Page     int
	Appendix string
}

type baseMap map[string]string
type RequestMap map[string]interface{}
