// Email sender, listen on RabbitMQ queue named "email" and
// send the email via mailjet service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	go statusMSG()
	
	viper.AutomaticEnv()
	viper.SetDefault("SERVICE_NAME", "srv-email")
	configFile := flag.String("config", "./config.toml", "path of the config file")
	flag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Printf("cannot read config file: %v\nUse env instead\n", err)
	}
	
	// Config email sending
    mailjetClient := mailjet.NewMailjetClient(viper.GetString("MJ_APIKEY_PUBLIC"), viper.GetString("MJ_APIKEY_PRIVATE"))

	srvName := viper.GetString("SERVICE_NAME")
	from := viper.GetString("EMAIL_FROM")
	
	// RabbitMQ email channel
	conn, err := amqp.Dial(viper.GetString("RABBITMQ_URI"))
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

		messagesInfo := []mailjet.InfoMessagesV31 {
			mailjet.InfoMessagesV31{
			  From: &mailjet.RecipientV31{
				Email: from,
				Name: "Admin YIC",
			  },
			  To: &mailjet.RecipientsV31{
				mailjet.RecipientV31 {
				  Email: msg.Headers["To"].(string),
				  },
			  },
			  Subject: msg.Headers["Subject"].(string),
			  TextPart: "",
			  HTMLPart: string(msg.Body),
			},
		  }
		  messages := mailjet.MessagesV31{Info: messagesInfo }
		  res, err := mailjetClient.SendMailV31(&messages)
		  if err != nil {
			  log.Fatal(err)
		  }
		  fmt.Printf("Data: %+v\n", res)



		
		msg.Ack(false)
	}

	select {} // wait infinitly

}


func statusMSG() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "email service is alive")
	})

	http.ListenAndServe(":8080", nil)
}
