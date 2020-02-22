package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// SearchData struct
type SearchData struct {
	Search string   `json:"search"`
	Sites  []string `json:"sites"`
}

func searchWeb(pattern string, urls []string) []string {
	result := make([]string, 0, len(urls))
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		str := string(body)
		if strings.Contains(str, pattern) {
			result = append(result, url)
		}
	}
	return result
}

func searchPagesHandler(wr http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(wr, "Only POST supported", http.StatusBadRequest)
	}

	var data SearchData
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	}

	wr.Header().Add("Content-Type", "application/json")
	results := searchWeb(data.Search, data.Sites)

	jsonResults, err := json.Marshal(results)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	}
	wr.Write(jsonResults)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", searchPagesHandler)

	port := 8080
	fmt.Printf("Starting web server at: %v\n", port)
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
