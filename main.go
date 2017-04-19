package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":7070", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "aloha world!")
}
