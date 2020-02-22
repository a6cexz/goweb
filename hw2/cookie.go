package main

import (
	"fmt"
	"log"
	"net/http"
)

func firstHandler(wr http.ResponseWriter, request *http.Request) {
	cookie := &http.Cookie{
		Name:  "MyTestCookie",
		Value: "MyTestCookieValue",
	}
	http.SetCookie(wr, cookie)
	fmt.Fprintln(wr, "Done")
}

func secondHandler(wr http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("MyTestCookie")
	if err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprintln(wr, cookie.Value)
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/first", firstHandler)
	router.HandleFunc("/second", secondHandler)

	port := 8080
	fmt.Printf("Starting web server at: %v\n", port)
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
