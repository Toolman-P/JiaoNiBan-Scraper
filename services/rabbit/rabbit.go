package rabbit

import (
	"JiaoNiBan-data/databases"
	"JiaoNiBan-data/scrapers/base"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Push(shrefs *[]base.ScraperHref) {
	conn, err := amqp.Dial(amqp_addr)

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}

	defer func() {
		ch.Close()
		conn.Close()
	}()

	q, err := ch.QueueDeclare(queueid, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}
	for _, sh := range *shrefs {
		if f, _ := databases.CheckHrefExists(sh.Author, sh.Hash); !f {
			data, err := json.Marshal(sh)
			if err != nil {
				panic(err)
			}
			ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json; charset=utf-8",
					Body:        data,
				},
			)
		}
	}
}
