package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s guest:guest@localhost:5672 email@address.com \"subject with space then the body Ctrl+D to send the email\"", os.Args[0])
	}
	rabbitmqHost := os.Args[1]
	email := os.Args[2]
	subject := os.Args[3]
	body, err := bufio.NewReader(os.Stdin).ReadString(0)
	if err != io.EOF {
		failOnError(err, "Failed to parse email content")
	}

	conn, err := amqp.Dial("amqps://" + rabbitmqHost + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

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

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/html",
			Headers:       amqp.Table{"To": email, "Subject": subject},
			Body:          []byte(body),
			AppId:         "cli_send_email",
			CorrelationId: correlationID(),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("Sent To <%s> Subject \"%s\" Content:\n%s\n", email, subject, body)
}

func correlationID() string {
	c := make([]byte, 8)
	rand.Read(c)
	return base64.StdEncoding.EncodeToString(c)
}
