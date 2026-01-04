package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/slapxxi/snippetbox/internal/models"
)

type config struct {
	addr      string
	dsn       string
	staticDir string
}

type application struct {
	snippets       *models.SnippetModel
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	db             *sql.DB
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:passworD1234$@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "HTTP network address")
	flag.Parse()

	app := NewApplication(cfg.dsn)
	defer app.db.Close()

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLogger,
		Handler:  app.routes(),
	}

	app.infoLogger.Printf("starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	app.errorLogger.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
