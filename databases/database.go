package databases

import (
	"JiaoNiBan-data/scrapers/base"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbs struct {
	rdb *redis.Client
	mdb *mongo.Client
}

var data dbs

func Init() error {
	data.rdb = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: "",
		DB:       0,
	})
	_, err := data.rdb.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	data.mdb, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_addr))
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	err := data.rdb.Close()
	if err != nil {
		log.Fatal("Something happened when closing redis.")
		defer data.mdb.Disconnect(context.TODO())
		return err
	}
	err = data.mdb.Disconnect(context.TODO())

	if err != nil {
		log.Fatal("Something wrong happened when closing mongo.")
		return err
	}

	return nil
}

func CheckConnection() bool {
	_, err := data.rdb.Ping(context.TODO()).Result()
	if err != nil {
		return false
	}
	return err == nil
}

func CheckHrefExists(cat string, hash string) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}

	return data.rdb.SIsMember(context.TODO(), cat, hash).Result()
}

func AddHref(cat string, hash string) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}
	f, err := data.rdb.SAdd(context.TODO(), cat, hash).Result()
	return f == 1, err
}

func GetVersion(cat string) string {
	v := fmt.Sprintf("%s.sha256", cat)
	if i, _ := data.rdb.Exists(context.TODO(), v).Result(); i == 1 {
		r, _ := data.rdb.Get(context.TODO(), v).Result()
		return r
	}
	return "X"
}

func SetVersion(cat string, ver string) {
	v := fmt.Sprintf("%s.sha256", cat)
	data.rdb.Set(context.TODO(), v, ver, 0)
}

func AddPage(sc base.ScraperContent) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}
	c := data.mdb.Database("offices").Collection(sc.Author)
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
