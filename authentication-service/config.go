package main

import (
	"authentication-service/data"
	"database/sql"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}
