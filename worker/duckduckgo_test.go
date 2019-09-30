package main

import (
	"fmt"
	"testing"
)

func Test_duckduckgoSearch(t *testing.T) {
	googleResults, err := duckduckgoSearch("golang")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Println("Successfully got the result.")
		fmt.Println(googleResults)
	}
}

func Test_wikipediaSearch(t *testing.T) {
	result := getWikipediaResult("golang")
	if len(result) == 0 {
		t.Fail()
	}
}
