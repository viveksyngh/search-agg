package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/viveksyngh/search-api/rabbitmq"
)

const (
	//InProgressStatus in progress status
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

			var lastInsertID int64
			err = s.DB.QueryRow(`INSERT INTO searchquery(status, query, created_on) VALUES ($1, $2, $3) RETURNING id`,
				InProgressStatus, searchQuery.Query, time.Now()).Scan(&lastInsertID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
				return
			}

			queryMessage := QueryMessage{ID: lastInsertID, Query: searchQuery.Query}
			messageBytes, err := json.Marshal(queryMessage)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
				return
			}
			rabbitmq.PublishMessage(s.Queue, QueueName, messageBytes)

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

func (s *Server) handleSearchResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var searchResult SearchResult
		queryParam, ok := r.URL.Query()["queryId"]
		if !ok || len(queryParam[0]) < 1 {
			http.Error(w, "query paramater queryId is missing", http.StatusBadRequest)

		}
		queryID := queryParam[0]

		var (
			id        int
			status    string
			query     string
			createdOn time.Time
			title     string
			url       string
			text      string
		)

		if r.Method == http.MethodGet {
			rows, err := s.DB.Query(`SELECT id, status, query, created_on FROM searchquery WHERE id = $1`, queryID)
			defer rows.Close()

			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}

			if rows.Next() {
				if err = rows.Scan(&id, &status, &query, &createdOn); err != nil {
					http.Error(w, "Something went wrong", http.StatusInternalServerError)
					return
				}

				searchResult = SearchResult{
					ID:        id,
					Status:    status,
					Query:     query,
					CreatedOn: createdOn,
				}

			} else {
				http.Error(w, "query with given queryId not found", http.StatusNotFound)
				return
			}

			rows, err = s.DB.Query(`SELECT id, title, url FROM googlesearchresult WHERE  searchquery_id = $1`, queryID)
			if err != nil {
				http.Error(w, "Something went wrong "+err.Error(), http.StatusInternalServerError)
				return
			}

			var googleResults []Result
			for rows.Next() {
				if err = rows.Scan(&id, &title, &url); err != nil {
					http.Error(w, "Something went wrong "+err.Error(), http.StatusInternalServerError)
					return
				}

				result := Result{
					ID:    id,
					Title: title,
					URL:   url,
				}
				googleResults = append(googleResults, result)
			}
			searchResult.Google = googleResults

			var duckduckResults []Result
			rows, err = s.DB.Query(`SELECT id, title, url FROM duckduckgosearchresult WHERE  searchquery_id = $1`, queryID)
			if err != nil {
				http.Error(w, "Something went wrong "+err.Error(), http.StatusInternalServerError)
				return
			}
			for rows.Next() {
				if err = rows.Scan(&id, &title, &url); err != nil {
					http.Error(w, "Something went wrong : "+err.Error(), http.StatusInternalServerError)
					return
				}

				result := Result{
					ID:    id,
					Title: title,
					URL:   url,
				}
				duckduckResults = append(duckduckResults, result)
			}
			searchResult.DuckDuckGo = duckduckResults

			var wikipediaResult WikipediaResult
			rows, err = s.DB.Query(`SELECT id, result FROM wikipediasearchresult WHERE  searchquery_id = $1`, queryID)
			if err != nil {
				http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if rows.Next() {
				if err = rows.Scan(&id, &text); err != nil {
					http.Error(w, "Something went wrong "+err.Error(), http.StatusInternalServerError)
					return
				}

				wikipediaResult = WikipediaResult{
					ID:   id,
					Text: text,
				}

			}
			searchResult.Wikipedia = wikipediaResult

			jsonData, err := json.Marshal(searchResult)
			if err != nil {
				http.Error(w, "Something went wrong "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
