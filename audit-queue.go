package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	rabbitConnection = "amqp://guest:guest@audit-mq:5672/"
)

func initQueues() {

	var err error
	var rmqConn *amqp.Connection
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		rmqConn, err = amqp.Dial("amqp://guest:guest@audit-mq:5672/")
		if err == nil {
			break
		}
		log.Println(err)
	}

	if err != nil {
		failOnError(err, "Failed to rmqConnect to RabbitMQ")
	}

	go receiveError(rmqConn)
	go receiveUser(rmqConn)
	go receiveTransaction(rmqConn)
	go receiveQuote(rmqConn)
}

func receiveError(c *amqp.Connection) {

	ch, err := c.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"error_queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			a := errorEvent{}
			err := json.Unmarshal(d.Body, &a)

			if err == nil {
				errorEventHandler(a)
			}
		}
	}()

	<-forever
}

func receiveUser(c *amqp.Connection) {

	ch, err := c.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"user_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			a := userCommand{}
			err := json.Unmarshal(d.Body, &a)

			if err == nil {
				userCommandHandler(a)

			}
		}
	}()

	<-forever
}

func receiveTransaction(c *amqp.Connection) {

	ch, err := c.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"transaction_queue", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			a := accountTransaction{}
			err := json.Unmarshal(d.Body, &a)

			if err == nil {
				accountTransactionHandler(a)
			}
		}
	}()

	<-forever
}

func receiveQuote(c *amqp.Connection) {

	ch, err := c.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"quote_queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			a := quoteServer{}
			err := json.Unmarshal(d.Body, &a)

			fmt.Println("quote mq got", err)

			if err == nil {
				fmt.Println("quote mq got")
				quoteServerHandler(a)
			}
		}
	}()

	<-forever
}
