package main

import (
	"RESTapi/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

func main() {
	const op = "storage.Connection()"

	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Формирование строки подключения
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	// Открытие соединения с базой данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	// Проверка на успешное подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("%s %v", op, err)
	}

	// Отложенное закрытие соединения
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("%s: error closing database connection: %v", op, err)
		}
	}()

	// Пример использования базы данных
	fmt.Println("Database connection established successfully")
}
