package main

import (
	"../../pkg/models/mysql"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//
//
//
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippetsStore *mysql.SnippetsStore
	templateCache map[string]*template.Template
}

//
//
//
func main() {

	// Startup arguments (flags)
	addr := flag.String("addr", ":4000", "HTTP Listening Address")
	dsn := flag.String("dsn", "snippetbox:box@/snippetbox?parseTime=true", "MySQL Data Source Name")
	flag.Parse()

	// Loggers init.
	infoLog := log.New(os.Stdout, "INFO  ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	// Database init.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Template cache init.
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatalf("Template cache init error: %s", err.Error())
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippetsStore: &mysql.SnippetsStore{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s ...\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

//
//
//
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
