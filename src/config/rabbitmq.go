package config

import "github.com/streadway/amqp"

func ConnectRabbit() (conn *amqp.Connection) {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	HandlerError(err, "can't connect to RabbitMQ")
	return
}

func CreateChannel(conn *amqp.Connection) (channel *amqp.Channel) {
	var err error
	channel, err = conn.Channel()
	HandlerError(err, "can't create a channel")
	err = channel.Qos(0, 0, true)
	HandlerError(err, "Could not configure QoS")
	return
}

func DeclareQueue(channel *amqp.Channel) (queue amqp.Queue) {
	var err error
	HandlerError(err, "Could not declare exchange")
	queue, err = channel.QueueDeclare("file", false, false, false, false, nil)
	return
}

func CloseRabbit(conn *amqp.Connection, channel *amqp.Channel) {
	conn.Close()
	channel.Close()
}
