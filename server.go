package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/viveksyngh/search-api/db"
)

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

func (s *Server) handleSearchQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var searchQuery SearchQuery

		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can not read body", http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(body, &searchQuery)
			if err != nil {
				http.Error(w, "malformed request body", http.StatusBadRequest)
				return
			}
			fmt.Println(searchQuery.Query)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
