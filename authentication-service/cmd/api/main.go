package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8083"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

var counts int64

func main() {
	log.Println("Starting authentication service...")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

/* Creates the DB connection and return it */
func openDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

/* Tries to connect to the DB and backs off every four seconds */
func connectToDB() *sql.DB {
	connectionString := os.Getenv("DSN")

	/* Stay in this loop until we connect to the DB */
	for {
		connection, err := openDB(connectionString)
		if err != nil {
			log.Println("Postgres is not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds")
		time.Sleep(4 * time.Second)

		continue

	}
}
