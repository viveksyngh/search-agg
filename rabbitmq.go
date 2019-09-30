package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

//Connection get an amqp connection
func Connection(queueName string) (*amqp.Connection, error) {
	//Connecting to rabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect rabbitMQ")

	//Opening channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	//Creating a Queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare Queue")

	return conn, err
}

//PublishMessage publish message to the queue
func PublishMessage(conn *amqp.Connection, queueName string, message []byte) {
	//Opening channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	err = ch.Publish(
		"",        // exchange
		queueName, // routing queue
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(message),
		})

	failOnError(err, "Failed to publish message")
}
