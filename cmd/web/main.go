package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
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
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:passworD1234$@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "HTTP network address")
	flag.Parse()

	app := &application{
		infoLogger:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	db, err := openDB(cfg.dsn)
	if err != nil {
		app.errorLogger.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLogger.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app.snippets = &models.SnippetModel{DB: db}
	app.templateCache = templateCache
	app.formDecoder = formDecoder
	app.sessionManager = sessionManager

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLogger,
		Handler:  app.routes(),
	}

	app.infoLogger.Printf("starting server on %s", cfg.addr)
	err = srv.ListenAndServe()
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
