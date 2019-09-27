package main

import (
	"database/sql"
	"net/http"

	"github.com/viveksyngh/search-api/db"
)

//CreateNewServer create and return a new server instance
func CreateNewServer() (*Server, error) {
	db, err := db.Connection()
	if err != nil {
		return &Server{}, err
	}

	return &Server{
		DB:     db,
		Router: http.NewServeMux(),
	}, nil
}

//Start start the HTTP Server
func (s *Server) Start() {
	http.ListenAndServe(":8000", s.Router)
}

//Server application server struct
type Server struct {
	DB     *sql.DB
	Router *http.ServeMux
}

//SearchQuery search query struct
type SearchQuery struct {
	Query string `json:"query"`
}

func (s *Server) routes() {
	s.Router.HandleFunc("/search", s.handleSearchQuery())
}
