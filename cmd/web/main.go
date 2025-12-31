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

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "HTTP network address")
	flag.Parse()

	app := &application{
		infoLogger:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := app.routes()

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLogger,
		Handler:  mux,
	}

	app.infoLogger.Printf("starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	app.errorLogger.Fatal(err)
}
