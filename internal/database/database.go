package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func CreateDatabase() (*sql.DB, error) {
	// Получение пути к базе данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db" // Путь по умолчанию
	}

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(LENGTH(repeat) <= 128)
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return nil, err
		}
		log.Println("База данных и таблица успешно созданы.")
	}

	return db, nil
}
