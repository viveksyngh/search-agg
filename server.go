package main

import (
	"database/sql"
	"net/http"

	"github.com/streadway/amqp"
	"github.com/viveksyngh/search-api/db"
	"github.com/viveksyngh/search-api/rabbitmq"
)

const (
	QueueName = "task_queue"
)

//CreateNewServer create and return a new server instance
func CreateNewServer() (*Server, error) {
	db, err := db.Connection()
	if err != nil {
		return &Server{}, err
	}

	queue, err := rabbitmq.Connection(QueueName)
	if err != nil {
		return &Server{}, err
	}
	return &Server{
		DB:     db,
		Router: http.NewServeMux(),
		Queue:  queue,
	}, nil
}

//Start start the HTTP Server
func (s *Server) Start() {
	defer s.DB.Close()
	defer s.Queue.Close()
	http.ListenAndServe(":8000", s.Router)
}

//Server application server struct
type Server struct {
	DB     *sql.DB
	Router *http.ServeMux
	Queue  *amqp.Connection
}

func (s *Server) routes() {
	s.Router.HandleFunc("/search", s.handleSearchQuery())
	s.Router.HandleFunc("/recent", s.handleRecentQueries())
	s.Router.HandleFunc("/search-result", s.handleSearchResult())
}
