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

type Result struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type WikipediaResult struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type SearchResult struct {
	ID         int             `json:"id"`
	Status     string          `json:"status"`
	Query      string          `json:"query"`
	CreatedOn  time.Time       `json:"created_on"`
	Google     []Result        `json:"google"`
	DuckDuckGo []Result        `json:"duckduckgo"`
	Wikipedia  WikipediaResult `json:"wikipedia"`
}
