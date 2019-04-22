package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func helloServrer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside HelloServer handler")
	fmt.Fprintf(w, "<h1>Hellow!"+req.URL.Path[1:]+"</h1>")
}

func testServrer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside HelloServer handler")
	fmt.Fprintf(w, "<h1>Test!"+req.URL.Path[1:]+"</h1>")
}

type empty struct {
}

func (e empty) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Inside HelloServer handler")
	paths := strings.Split(req.URL.Path, "/")
	fmt.Fprintf(w, "<h1>Hello "+paths[len(paths)-1]+"</h1>")
}

func main() {
	http.HandleFunc("/test", testServrer)
	// http.HandleFunc("/", helloServrer)
	var e empty
	err := http.ListenAndServe("localhost:8080", e)
	if err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
