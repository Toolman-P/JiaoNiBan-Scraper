package databases

import (
	"JiaoNiBan-data/scrapers/base"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	mdb, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_addr))
	if err != nil {
		return err
	}

	return nil
}

func Close() {

	var err error
	if rdb != nil {
		if err = rdb.Close(); err != nil {
			panic(err)
		}
	}

	if mdb != nil {
		if err = mdb.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func CheckConnection() bool {
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		return false
	}
	return err == nil
}

func CheckHrefExists(opt string, hash string) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}

	return rdb.SIsMember(context.TODO(), opt, hash).Result()
}

func AddHref(opt string, hash string) (bool, error) {
	if !CheckConnection() {
		return false, errors.New("connection failed")
	}
	f, err := rdb.SAdd(context.TODO(), opt, hash).Result()
	return f == 1, err
}

func GetVersion(opt string) string {
	v := fmt.Sprintf("%s.sha256", opt)
	if i, _ := rdb.Exists(context.TODO(), v).Result(); i == 1 {
		r, _ := rdb.Get(context.TODO(), v).Result()
		return r
	}
	return "X"
}

func GetLatestPage(opt string) int {
	v := fmt.Sprintf("%s.latest", opt)
	if i, _ := rdb.Exists(context.TODO(), v).Result(); i == 1 {
		r, _ := rdb.Get(context.TODO(), v).Int()
		return r
	}
	return -1
}

func GetPageSum(opt string) int {
	v := fmt.Sprintf("%s.sum", opt)
	if i, _ := rdb.Exists(context.TODO(), v).Result(); i == 1 {
		r, _ := rdb.Get(rdb.Context(), v).Int()
		return r
	}
	return -1
}

func SetVersion(opt string, ver string) {
	v := fmt.Sprintf("%s.sha256", opt)
	rdb.Set(context.TODO(), v, ver, 0)
}

func SetLatestPage(opt string, page int) {
	v := fmt.Sprintf("%s.latest", opt)
	rdb.Set(context.TODO(), v, page, 0)
}

func SetPageSum(opt string, page int) {
	v := fmt.Sprintf("%s.sum", opt)
	rdb.Set(context.TODO(), v, page, 0)
}

func AddDesc(opt string, sc *base.ScraperContent) error {
	if !CheckConnection() {
		return errors.New("connection failed")
	}
	c := mdb.Database(descs).Collection(opt)
	_, err := c.InsertOne(context.TODO(), bson.D{{"title", sc.Title},
		{"author", sc.Author},
		{"year", sc.Year},
		{"month", sc.Month},
		{"day", sc.Day},
		{"page", sc.Page},
		{"description", sc.Desc},
		{"sha256", sc.Hash}})

	if err != nil {
		return err
	}
	return nil
}

func AddContent(opt string, sc *base.ScraperContent) error {
	if !CheckConnection() {
		return errors.New("connection failed")
	}

	c := mdb.Database(contents).Collection(opt)
	_, err := c.InsertOne(context.TODO(), bson.D{{"sha256", sc.Hash}, {"body", sc.Body}, {"appendix", sc.Appendix}})
	if err != nil {
		return err
	}
	return nil
}
