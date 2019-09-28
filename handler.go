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
