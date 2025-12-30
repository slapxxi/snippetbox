package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "HTTP network address")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(cfg.staticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets", snippetView)
	mux.HandleFunc("/snippets/new", snippetCreate)

	srv := http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLogger,
		Handler:  mux,
	}

	infoLogger.Printf("starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	errorLogger.Fatal(err)
}
