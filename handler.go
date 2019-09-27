package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
			fmt.Println(searchQuery.Query)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
