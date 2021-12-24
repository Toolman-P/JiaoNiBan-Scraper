package databases

import (
	"os"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

var rdb *redis.Client
var mdb *mongo.Client

var redis_addr = os.Getenv("REDIS_ADDR")
var mongo_addr = os.Getenv("MONGO_ADDR")
