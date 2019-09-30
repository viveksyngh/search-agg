package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func buildDuckDuckGoURL(query string) string {
	baseURL := "https://duckduckgo.com/html"
	url, _ := url.Parse(baseURL)
	queryParams := url.Query()
	queryParams.Set("q", query)
	queryParams.Set("s", "0")
	url.RawQuery = queryParams.Encode()
	return url.String()
}

func duckduckGo(keyword string) (*http.Response, error) {
	baseClient := &http.Client{}
	url := buildDuckDuckGoURL(keyword)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	res, err := baseClient.Do(req)

	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func duckduckgoResultParser(response *http.Response) ([]SearchResult, error) {
	// defer response.Body.Close()
	// responseText, err := ioutil.ReadAll(response.Body)
	// fmt.Println(string(responseText))
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []SearchResult{}
	sel := doc.Find("div#links .links_main")
	fmt.Println(sel)
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		title := linkTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" {
			result := SearchResult{
				Title: title,
				URL:   link,
			}
			results = append(results, result)
		}
	}
	return results, err
}

func duckduckgoSearch(query string) ([]SearchResult, error) {
	results := []SearchResult{}
	res, err := duckduckGo(query)
	if err != nil {
		return results, err
	}

	return duckduckgoResultParser(res)
}
