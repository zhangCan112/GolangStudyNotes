package main

import (
	"GolangStudyNotes/negroni"
	"fmt"
	"log"
	"net/http"
)

// middleware1 用来返回标题部分的内容
func middleware1(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "<h1>标题：你好！！！</h1>")
}

func middleware2(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "<p>这是正文部分中间件产生的内容</p>")
}

func middleware3(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "<p>这里结尾部分的内容</p>")
}

func main() {
	handler := negroni.New(
		negroni.WrapFunc(middleware1),
		negroni.WrapFunc(middleware2),
		negroni.WrapFunc(middleware3),
	)

	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err.Error())
	}
}
