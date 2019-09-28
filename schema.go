package main

import "time"

//SearchItem search query item
type SearchItem struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Query     string    `json:"query"`
	CreatedOn time.Time `json:"created_on"`
}

type QueryMessage struct {
	ID    int64  `json:"id"`
	Query string `json:"query"`
}

//SearchQuery search query struct
type SearchQuery struct {
	Query string `json:"query"`
}
