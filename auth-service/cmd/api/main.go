package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nico2220/auth-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)


var (
	count = 0
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
		// maxOpenConns int
		// maxIdleConns int
		// maxIdleTime time.Duration
	}
}

type application struct {
	config config
	Models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "api environment")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "postgres db dsn")

	flag.Parse()

	log.Printf("starting authentication service on port %d \n", cfg.port)

	conn := connectToDb(cfg)
	if conn == nil {
		log.Panic()
	}

	app := application{
		config: cfg,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDb(cfg config) *sql.DB {

	for {

		conn, err := openDb(cfg.db.dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			count++
		} else {
			log.Println("Connected to Postgres...")
			return conn
		}

		if count >= 10 {
			log.Println(err)
			return nil
		}

		time.Sleep(time.Second * 1)

	}

}
