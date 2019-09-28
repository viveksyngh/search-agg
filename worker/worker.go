package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {

	//Connecting to rabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect rabbitMQ")
	defer conn.Close()

	//Opening channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	//Creating a Queue
	queue, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // autoDelete
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	failOnError(err, "Failed to declare Queue")

	//Setting prefetch count for fair load balancing
	err = ch.Qos(
		1,     //Prefetch count
		0,     //Prefetch size
		false, //globale
	)
	failOnError(err, "Failed to set prefetch couny")

	messages, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to publish message")

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var searchQuery SearchQuery
			log.Printf("Received a message: %s", d.Body)
			err := json.Unmarshal(d.Body, &searchQuery)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				search(searchQuery)
				log.Printf("Done!")
			}
			d.Ack(false)
		}
	}()

	log.Printf("[*] Wating for messages. To exit press CTRL+C")
	<-forever
}

type SearchResult struct {
	Title string
	Link  string
}

type SearchQuery struct {
	ID    int64  `json:"id"`
	Query string `json:"query"`
}

func search(searchQuery SearchQuery) {

	googleChan := make(chan []SearchResult)
	duckduckgoChan := make(chan []SearchResult)
	wikipediaChan := make(chan string)

	go func() {
		googleChan <- googleSearch(searchQuery.Query)
	}()

	go func() {
		duckduckgoChan <- duckDuckGoSearch(searchQuery.Query)
	}()

	go func() {
		wikipediaChan <- wikipediaSearch(searchQuery.Query)
	}()

	timeout := time.After(80 * time.Millisecond)

	select {
	case results := <-googleChan:
		//Insert the results in the database
		fmt.Println(results)
	case results := <-duckduckgoChan:
		//Insert the results in the database
		fmt.Println(results)
	case wikipediaResult := <-wikipediaChan:
		//Insert the result in the database
		fmt.Println(wikipediaResult)
	case <-timeout:
		fmt.Println("Search Timed Out")
		return
	}

}

func duckDuckGoSearch(query string) []SearchResult {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return []SearchResult{
		{
			Title: "Test",
			Link:  "www.duckduckgo.com",
		},
	}
}

func googleSearch(query string) []SearchResult {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return []SearchResult{
		{
			Title: "Test Search Result",
			Link:  "www.google.com",
		},
	}
}

func wikipediaSearch(query string) string {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return "Random test search result"
}
