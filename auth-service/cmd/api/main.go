package main

import (
	"database/sql"
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

const (
	Port = "8080"
)

var (
	count = 0
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	log.Printf("starting authentication service on port %s\n", Port)

	conn := connectToDb()
	if conn == nil {
		log.Panic()
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", Port),
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

func connectToDb() *sql.DB {
	dsn := os.Getenv("DSN")

	for {

		conn, err := openDb(dsn)
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
