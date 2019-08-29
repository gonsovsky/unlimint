package consumer

import (
	"../shared"
	"../storage"
	"github.com/streadway/amqp"
	"log"
)

//Consumer - subsribes for events
type Consumer struct {
	No     int
	Config *shared.Cfg
	Db     storage.IRepository
}

//Subscribe - Let's launch web server
func (consumer *Consumer) Subscribe() {
	conn, err := amqp.Dial(consumer.Config.Amqp.URL)
	failOnReceive(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	ch.Confirm(false)
	failOnReceive(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		consumer.Config.Amqp.Queue, // name
		true,                       // durable
		false,                      // delete when usused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnReceive(err, "Failed to declare a queue")

	hits, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnReceive(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range hits {
			hit := shared.GoogleHit{}
			hit.FromJSON(d.Body)
			consumer.doWork(hit)
			d.Ack(false)
		}
	}()

	<-forever
}

func (consumer *Consumer) doWork(hit shared.GoogleHit) {
	consumer.Db.Post(hit)
}

func failOnReceive(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
