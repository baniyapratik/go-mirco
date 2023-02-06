package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"authentication-service/data"
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = "8082"

var counts int64

func main() {
	log.Println("Starting authentication-service")
	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to postgres!")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    PORT,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
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

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet!")
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
