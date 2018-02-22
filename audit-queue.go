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

var (
	dumplogAudit = &userCommand{}

	emptiedQueues = make(chan int)
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
	go dumpLogReady(emptiedQueues)
}

func dumpLogReady(msg <-chan int) {
	var queuesFinished = 0
	for range msg {
		queuesFinished++
		fmt.Println("WOW dumplog COMMMIN UP", queuesFinished)
		if queuesFinished == 4 {
			//fmt.Println(dumplogAudit, "fdnjafndjajkflnjalfndajknfdljandl")
			userCommandHandler(*dumplogAudit) //log the command to dump
			dumpLogCommand()                  //big dump
		}
	}

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

			if err != nil {
				fmt.Println("marshal error")
			}
			if a.Username == "DUMPLOG" {
				emptiedQueues <- 1
			} else if err == nil {
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

			if err != nil {
				fmt.Println("marshal error")
			}

			if a.Command == "DUMPLOG" {
				emptiedQueues <- 1
				dumplogAudit = &a
				broadcastDumplog(c)
				fmt.Println(dumplogAudit, "Fdfdafdada")
			} else if err == nil {
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

			if err != nil {
				fmt.Println("marshal error")
			}
			if a.Username == "DUMPLOG" {
				emptiedQueues <- 1
			} else if err == nil {
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
			if err != nil {
				fmt.Println("marshal error")
			}
			if a.Username == "DUMPLOG" {
				emptiedQueues <- 1
			} else if err == nil {
				quoteServerHandler(a)
			}
		}
	}()

	<-forever
}

func broadcastDumplog(c *amqp.Connection) {

	m := errorEvent{Username: "DUMPLOG"}
	body, merr := json.Marshal(m)

	if merr != nil {
		fmt.Println("marshal error")
	}

	ch, err := c.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Publish(
		"",            // exchange
		"error_queue", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
		})
	failOnError(err, "Failed to publish error mq")

	a := accountTransaction{Username: "DUMPLOG"}
	body, merr = json.Marshal(a)

	err = ch.Publish(
		"",                  // exchange
		"transaction_queue", // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
		})
	failOnError(err, "Failed to publish error mq")

	q := quoteServer{Username: "DUMPLOG"}
	body, merr = json.Marshal(q)

	err = ch.Publish(
		"",            // exchange
		"quote_queue", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
		})
	failOnError(err, "Failed to publish error mq")

}
