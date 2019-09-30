package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getWikipediaURL(query string) (string, error) {
	googleResults, err := GoogleScrape(query, "com", "en")
	if err != nil {
		fmt.Println("Can not get results from google: ", err.Error())
		return "", err
	}

	for _, result := range googleResults {
		if strings.Contains(result.ResultURL, "wikipedia.org") {
			return result.ResultURL, nil
		}
	}

	return "", nil
}

func wikipediaResultParser(response *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return "", err
	}
	result := ""
	sel := doc.Find("div#mw-content-text p")
	index := 0
	for i := range sel.Nodes {
		if index > 1 {
			break
		}
		item := sel.Eq(i)
		result = result + item.Text() + "\n"
		index += 1

	}
	return result, nil
}

func getWikipediaResult(query string) string {
	url, err := getWikipediaURL(query)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	res, err := googleRequest(url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	result, err := wikipediaResultParser(res)
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
