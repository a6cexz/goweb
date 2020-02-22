package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func searchWeb(pattern string, urls []string) []string {
	result := []string{}
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

func main() {
	urls := []string{"https://golang.org/pkg/bufio/", "https://golang.org/pkg"}
	r := searchWeb("ScanLines", urls)
	for _, s := range r {
		fmt.Println(s)
	}
}
