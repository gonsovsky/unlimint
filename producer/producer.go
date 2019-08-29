package producer

import (
	"../shared"
	"github.com/streadway/amqp"
	"log"
)

type Producer struct {
	config *shared.Cfg
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewProducer(config *shared.Cfg) *Producer {
	x := &Producer{config: config}
	x.Open()
	return x;
}

func failOnSend(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//Open channel to RabbitMQ
func (sender *Producer) Open() {
	conn, err := amqp.Dial(sender.config.Amqp.URL)
	failOnSend(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	ch.Confirm(false)
	failOnSend(err, "Failed to open a channel")
	ch.QueueDelete(sender.config.Amqp.Queue,false,false, false)

	queue, err := ch.QueueDeclare(
		sender.config.Amqp.Queue, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	failOnSend(err, "Failed to declare a queue")

	sender.channel = ch
	sender.queue = &queue
}

//Publish - Send Message to RabbitMQ
func (sender *Producer) Publish(hit shared.GoogleHit) {
	ch := sender.channel
	queue := sender.queue

	err := ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate

		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         hit.ToJSON(),
			DeliveryMode: 2, //persistent
		})
	failOnSend(err, "Failed to publish a message")

	if err != nil {
		log.Printf("Failed to publish message to queue %s: %s", queue.Name, err)
		return
	}
}


