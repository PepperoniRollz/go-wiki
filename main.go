package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var (
	maxWorkers = 10
	sem        = make(chan struct{}, maxWorkers)
)

func main() {

	data, err := os.ReadFile("graph.json")
	if err != nil {
		//create graph a map cotaining all wikipedia articles
		articles := make(map[string][]string)
		visited := make(map[string]bool)
		startingUri := "https://en.wikipedia.org/wiki/Special:AllPages"
		visited[startingUri] = true
		crawlAllPages(startingUri, articles, visited)

		json, err := json.Marshal(articles)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.WriteFile("graph.json", json, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("graph saved to file successfully.")
	}

	//handle creating the adjacency list

	var articles map[string]json.RawMessage
	unmarshallingError := json.Unmarshal(data, &articles)

	if unmarshallingError != nil {
		log.Fatal(err)
	}
	fmt.Println(len(articles))

	adjList := make(map[string]map[string]bool)
	var wg sync.WaitGroup
	for k := range articles {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			worker("https://en.wikipedia.org"+url, adjList)

		}(k)
	}
	wg.Wait()

}

func crawlAllPages(url string, graph map[string][]string, visited map[string]bool) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var nextPage string

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var getArticles func(*html.Node)
	getArticles = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.HasPrefix(a.Val, "/wiki") && !strings.Contains(a.Val, ":") {

					_, ok := graph[a.Val]
					if !ok {
						graph[a.Val] = []string{}
					}
				}
				if a.Key == "href" && strings.Contains(a.Val, "Special:AllPages&from") && !strings.Contains(a.Val, "view_mobile") {
					nextPage = a.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getArticles(c)
		}
	}
	getArticles(doc)
	_, ok := visited[nextPage]
	if !ok {
		visited[nextPage] = true
		crawlAllPages("https://en.wikipedia.org"+nextPage, graph, visited)
	}
}

func getLinks(url string) (map[string]bool, error) {
	links := make(map[string]bool)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var link func(*html.Node)
	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.HasPrefix(a.Val, "/wiki") && !strings.Contains(a.Val, ":") {
					_, ok := links[a.Val]
					if !ok {
						links[a.Val] = true
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link(c)
		}
	}
	link(doc)
	return links, nil

}

func worker(url string, newMap map[string]map[string]bool) {
	sem <- struct{}{}
	defer func() { <-sem }()

	links, err := getLinks(url)
	if err != nil {
		newMap[url] = links
	}
	fmt.Println(url, len(links))

}
