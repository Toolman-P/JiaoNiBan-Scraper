package databases

import "os"

const (
	// redis_addr = "172.17.0.1:6379"
	// mongo_addr = "mongodb://172.17.0.1:27017/"
	descs    = "descs"
	contents = "contents"
)

var redis_addr = os.Getenv("REDIS_ADDR")
var mongo_addr = os.Getenv("MONGO_ADDR")
