package base

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36"
const DeanBaseURL = "http://www.jwc.sjtu.edu.cn"
const DeanFirstPage = "http://www.jwc.sjtu.edu.cn/index/mxxsdtz.htm"
const DeanPrefix = "http://www.jwc.sjtu.edu.cn/index/mxxsdtz"
const storagePath = "downloads/website"

type baseMap map[string]string
type RequestMap map[string]interface{}

var baseurlMap = baseMap{"dean": "http://www.jwc.sjtu.edu.cn"}
