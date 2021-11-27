package databases

import (
	"JiaoNiBan-data/scraper/base"
	"context"
	"errors"
	"log"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbs struct {
	Rdb *redis.Client
	Mdb *mongo.Client
}

var Data dbs

func Init() error {
	Data.Rdb = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: "",
		DB:       0,
	})
	_, err := Data.Rdb.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	Data.Mdb, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_addr))
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	err := Data.Rdb.Close()
	if err != nil {
		log.Fatal("Something happened when closing redis.")
		defer Data.Mdb.Disconnect(context.TODO())
		return err
	}
	err = Data.Mdb.Disconnect(context.TODO())

	if err != nil {
		log.Fatal("Something wrong happened when closing mongo.")
		return err
	}

	return nil
}

func CheckConnection() bool {
	_, err := Data.Rdb.Ping(context.TODO()).Result()
	if err != nil {
		return false
	}
	return err == nil
}

func CheckHrefExists(cat string, hash string) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}

	return Data.Rdb.SIsMember(context.TODO(), cat, hash).Result()
}

func AddHref(cat string, hash string) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}
	f, err := Data.Rdb.SAdd(context.TODO(), cat, hash).Result()
	return f == 1, err
}

func AddPage(sc *base.ScraperContent) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}
	c := Data.Mdb.Database("contents").Collection(sc.Author)
	_, err := c.InsertOne(context.TODO(), bson.D{{"title", sc.Title},
		{"author", sc.Author},
		{"date", sc.Date},
		{"description", sc.Description},
		{"sha256", sc.Hash},
		{"text", sc.Text}})

	if err != nil {
		return false, err
	}
	return true, nil
}
