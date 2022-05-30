package webservice

import (
	"file_processor/src/config"
	"file_processor/src/consumer"
	"file_processor/src/producer"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Incoming struct {
	NumberOfMessage int `json:"numberOfMessage"`
}

func TriggerMessages(c *gin.Context) {
	var wg sync.WaitGroup
	var incoming Incoming

	err := c.BindJSON(&incoming)
	config.HandlerError(err, "Can't bind JSON")

	total_routines := incoming.NumberOfMessage
	if total_routines > 1_000_000 {
		total_routines = 1_000_000
	}

	conn := config.ConnectRabbit()
	channel := config.CreateChannel(conn)

	message := map[string]string{
		"greetings": "Hello",
		"place":     "World",
	}

	producer.ProduceMessages(&wg, total_routines, channel, message)

	config.CloseRabbit(conn, channel)
	c.IndentedJSON(http.StatusOK, incoming)
}

func ReturnSomething(c *gin.Context) {
	var wg sync.WaitGroup

	conn := config.ConnectRabbit()
	channel := config.CreateChannel(conn)

	total_routines := consumer.TotalOfMessages(channel)
	wg.Add(total_routines)
	consumer.ConsumeMessages(channel, &wg)

	config.CloseRabbit(conn, channel)

	outcoming := map[string]int{
		"totalMessage": total_routines,
	}
	c.IndentedJSON(http.StatusOK, outcoming)
}
