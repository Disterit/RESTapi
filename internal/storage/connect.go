package storage

import (
	"RESTapi/internal/config"
	"database/sql"
	"fmt"
)

func Connection() *sql.DB {

	const op = "storage.Connection()"

	cfg := config.MustLoad()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Errorf("%s %w", op, err)
	}

	return db
}
