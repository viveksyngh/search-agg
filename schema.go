package main

import "time"

//SearchItem search query item
type SearchItem struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Query     string    `json:"query"`
	CreatedOn time.Time `json:"created_on"`
}
