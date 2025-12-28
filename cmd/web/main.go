package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/snippets", snippetView)
	http.HandleFunc("/snippets/new", snippetCreate)

	log.Print("starting server on :4000")
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}
