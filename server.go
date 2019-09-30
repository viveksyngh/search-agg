package main

import (
	"database/sql"
	"net/http"

	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"github.com/viveksyngh/search-agg/db"
)

const (
	//QueueName queue name for the tasks
	QueueName = "task_queue"
)

//CreateNewServer create and return a new server instance
func CreateNewServer() (*Server, error) {
	db, err := db.Connection()
	if err != nil {
		return &Server{}, err
	}

	queue, err := Connection(QueueName)
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
	handler := cors.Default().Handler(s.Router)
	http.ListenAndServe("0.0.0.0:8000", handler)
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
