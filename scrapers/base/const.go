package base

const (
	UserAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36"
	DeanBaseURL   = "http://www.jwc.sjtu.edu.cn"
	DeanFirstPage = "http://www.jwc.sjtu.edu.cn/index/mxxsdtz.htm"
	DeanPrefix    = "http://www.jwc.sjtu.edu.cn/index/mxxsdtz"
	storagePath   = "downloads/website"
)

type baseMap map[string]string
type RequestMap map[string]interface{}

var baseurlMap = baseMap{"dean": "http://www.jwc.sjtu.edu.cn"}
