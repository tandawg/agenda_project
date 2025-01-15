package main

import (
	"database/sql"
	"log"

	"github.com/tandawg/agenda_project/internal/database"
	"github.com/tandawg/agenda_project/internal/webserver"
)

var db *sql.DB // Глобальная переменная для подключения к базе данных

func main() {
	// Инициализация базы данных
	var err error
	db, err = database.CreateDatabase()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Запуск веб-сервера
	webserver.StartServer(db)
}
