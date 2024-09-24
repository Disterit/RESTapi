package storage

import (
	"RESTapi/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connection() *sql.DB {
	const op = "storage.Connection()"

	cfg := config.MustLoad()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("%s %v", op, err)
	}

	return db
}
