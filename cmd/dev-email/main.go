// Fake Email sender, listen on RabbitMQ queue named "email"  and print it on the standard output
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	srvName := "srv-email-print"

	// RabbitMQ email channel
	conn, err := amqp.Dial(os.Args[1])
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	go func() {
		log.Fatalf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
	}()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,  // queue
		srvName, // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Printf("%s ready, wait on amqp queue '%s'", srvName, q.Name)

	for msg := range msgs {
		to := msg.Headers["To"].(string)
		subject := msg.Headers["Subject"].(string)
		content := string(msg.Body)

		log.Printf("To: %s  CorrID=%s\nSubject: %s\n%s\n", to, msg.CorrelationId, subject, content)
		msg.Ack(false)
	}
}
