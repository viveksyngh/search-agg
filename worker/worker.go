package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
	"github.com/viveksyngh/search-agg/db"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {

	//Connecting to rabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
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

	dbConn, err := db.Connection()
	if err != nil {
		log.Fatal(err.Error())
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var searchQuery SearchQuery
			log.Printf("Received a message: %s", d.Body)
			err := json.Unmarshal(d.Body, &searchQuery)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				search(dbConn, searchQuery)
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
	URL   string
}

type SearchQuery struct {
	ID    int64  `json:"id"`
	Query string `json:"query"`
}

func search(db *sql.DB, searchQuery SearchQuery) {

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

	timeout := time.After(100000 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case results := <-googleChan:
			//Insert the results in the database
			// fmt.Println(results)
			for _, result := range results {
				var title string
				if len(result.Title) > 254 {
					title = result.Title[:254]
				} else {
					title = result.Title
				}
				_, err := db.Exec(`INSERT INTO googlesearchresult (title, searchquery_id, url) VALUES ($1, $2, $3)`,
					title, searchQuery.ID, result.URL)
				if err != nil {
					fmt.Println("Failed to insert value: ", err.Error())
				}
			}
		case results := <-duckduckgoChan:
			//Insert the results in the database
			for _, result := range results {
				var title string
				if len(result.Title) > 254 {
					title = result.Title[:254]
				} else {
					title = result.Title
				}
				_, err := db.Exec(`INSERT INTO duckduckgosearchresult (title, searchquery_id, url) VALUES ($1, $2, $3)`,
					title, searchQuery.ID, result.URL)
				if err != nil {
					fmt.Println("Failed to insert value: ", err.Error())
				}
			}
		case wikipediaResult := <-wikipediaChan:
			//Insert the result in the database
			fmt.Println(wikipediaResult)
			_, err := db.Exec(`INSERT INTO wikipediasearchresult (result, searchquery_id) VALUES ($1, $2)`,
				wikipediaResult, searchQuery.ID)
			if err != nil {
				fmt.Println("Failed to insert value: ", err.Error())
			}
		case <-timeout:
			fmt.Println("Search Timed Out")
			return
		}
	}

	_, err := db.Exec(`UPDATE searchquery SET status = $1 WHERE id = $2`, "Completed", searchQuery.ID)
	if err != nil {
		fmt.Println("Failed to update the status: ", err.Error())
	}
	return

}

func duckDuckGoSearch(query string) []SearchResult {
	// time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	results := []SearchResult{}
	results, err := duckduckgoSearch(query)
	if err != nil {
		fmt.Println("Can not get results from google: ", err.Error())
		return results
	}
	if len(results) > 10 {
		results = results[:10]
	}
	return results
}

func googleSearch(query string) []SearchResult {
	var results []SearchResult
	googleResults, err := GoogleScrape(query, "com", "en")
	if err != nil {
		fmt.Println("Can not get results from google: ", err.Error())
		return results
	}

	for _, googleResult := range googleResults {
		results = append(results, SearchResult{
			Title: googleResult.ResultDesc,
			URL:   googleResult.ResultURL})
	}
	if len(results) > 10 {
		results = results[:10]
	}
	return results
}

func wikipediaSearch(query string) string {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return "Random test search result"
}
