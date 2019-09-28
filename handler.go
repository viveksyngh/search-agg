package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	InProgressStatus = "In Progress"
)

//handleSearchQuery returns handler to handler search query request
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

			_, err = s.DB.Exec(`INSERT INTO searchquery(status, query, created_on) VALUES ($1, $2, $3)`,
				InProgressStatus, searchQuery.Query, time.Now())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
				return
			}

			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "Search query submitted.")

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *Server) handleRecentQueries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var searchItems []SearchItem
		var (
			id        int
			status    string
			query     string
			createdOn time.Time
		)

		if r.Method == http.MethodGet {
			rows, err := s.DB.Query(`SELECT id, status, query, created_on FROM searchquery ORDER BY created_on DESC LIMIT 10`)
			defer rows.Close()

			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				if err = rows.Scan(&id, &status, &query, &createdOn); err != nil {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}

				searchItems = append(searchItems, SearchItem{
					ID:        id,
					Status:    status,
					Query:     query,
					CreatedOn: createdOn})

			}

			jsonData, err := json.Marshal(searchItems)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
