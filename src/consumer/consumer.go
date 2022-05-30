package consumer

import (
	"encoding/json"
	"file_processor/src/config"
	"file_processor/src/repository"
	"sync"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConsumeMessages(channel *amqp.Channel, wg *sync.WaitGroup) {
	documentChannel := make(chan map[string]string)

	message, err := channel.Consume("file", "consumer-1", true, false, false, false, nil)
	config.HandlerError(err, "Could not register consumer")

	go consume(message, documentChannel, channel)

	watchdogQueue(channel)

	go saveOnMongo(documentChannel, wg)
	wg.Wait()
}

func consume(messages <-chan amqp.Delivery, documentChannel chan map[string]string, channel *amqp.Channel) {
	for message := range messages {
		extractMessage(message, documentChannel)
	}
}

func extractMessage(message amqp.Delivery, messageChannel chan map[string]string) {
	var body map[string]string

	err := json.Unmarshal(message.Body, &body)
	config.HandlerError(err, "Error decoding JSON")

	go func() {
		messageChannel <- body
	}()
}

func TotalOfMessages(channel *amqp.Channel) int {
	newQueue, err := channel.QueueInspect("file")
	config.HandlerError(err, "Can't inspect queue")
	return newQueue.Messages
}

func watchdogQueue(channel *amqp.Channel) {
	total := 1
	for total > 0 {
		total = TotalOfMessages(channel)
	}
	channel.Cancel("consumer-1", false)
}

func save(message map[string]string, collection *mongo.Collection) {
	doc := bson.D{
		primitive.E{Key: "Greeting", Value: string(message["greetings"])},
		primitive.E{Key: "Place", Value: string(message["place"])},
	}
	repository.InsertOne(doc, collection)
}

func saveOnMongo(documentChannel chan map[string]string, wg *sync.WaitGroup) {
	collection := repository.CollectionMongo()

	for document := range documentChannel {
		save(document, collection)
		wg.Done()
	}

	repository.CloseCollection()
}
