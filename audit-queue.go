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
	dumplogAudit    = &userCommand{}
	emptiedQueues   = make(chan int)
	errorBulk       = []errorEvent{}
	userBulk        = []userCommand{}
	transactionBulk = []accountTransaction{}
	quoteBulk       = []quoteServer{}
	bulkAmount      = 20
)

func initQueues() {

	var err error
	var rmqConn *amqp.Connection
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		rmqConn, err = amqp.Dial("amqp://guest:guest@192.168.1.143:5672/")
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
			queuesFinished = 0
			//fmt.Println(dumplogAudit, "fdnjafndjajkflnjalfndajknfdljandl")
			userBulk = append(userBulk, *dumplogAudit)

			//go func() {
			userCommandHandler(userBulk) //log the command to dump
			//<-semaphoreChan
			//}()

			//semaphoreChan <- struct{}{}
			//go func() {
			dumpLogCommand() //big dump
			//	<-semaphoreChan
			//}()

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
				go errorEventHandler(errorBulk)
				errorBulk = nil
			} else if err == nil {
				errorBulk = append(errorBulk, a)
				if len(errorBulk) > bulkAmount {
					go errorEventHandler(errorBulk)
					errorBulk = nil
				}

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
				go userCommandHandler(userBulk)
				userBulk = nil
			} else if err == nil {
				userBulk = append(userBulk, a)
				if len(userBulk) > bulkAmount {
					go userCommandHandler(userBulk)
					userBulk = nil
				}

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
				go accountTransactionHandler(transactionBulk)
				transactionBulk = nil
			} else if err == nil {

				transactionBulk = append(transactionBulk, a)
				if len(transactionBulk) > bulkAmount {
					go accountTransactionHandler(transactionBulk)
					transactionBulk = nil
				}

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
				go quoteServerHandler(quoteBulk)
				quoteBulk = nil
			} else if err == nil {
				quoteBulk = append(quoteBulk, a)
				if len(quoteBulk) > bulkAmount {
					go quoteServerHandler(quoteBulk)
					quoteBulk = nil
				}

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
