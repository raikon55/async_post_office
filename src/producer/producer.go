package producer

import (
	"encoding/json"
	"file_processor/src/config"
	"sync"

	"github.com/streadway/amqp"
)

func ProduceMessages(wg *sync.WaitGroup, total_routines int, channel *amqp.Channel, message map[string]string) {
	wg.Add(total_routines)
	for i := 0; i < total_routines; i++ {
		go produce(channel, message, wg)
	}
	wg.Wait()
}

func produce(channel *amqp.Channel, message map[string]string, wg *sync.WaitGroup) {
	body, err := json.Marshal(message)
	config.HandlerError(err, "Error encoding JSON")

	channel.Publish("file_exchange", "arq", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	go func() {
		wg.Done()
	}()
}
