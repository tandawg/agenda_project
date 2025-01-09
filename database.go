package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB // Глобальная переменная для базы данных

func init() {
	DB = createDatabase()
}

func createDatabase() *sql.DB {
    // Получение пути к базе данных из переменной окружения
    dbFile := os.Getenv("TODO_DBFILE")
    if dbFile == "" {
        dbFile = "scheduler.db" // Путь по умолчанию
    }

    _, err := os.Stat(dbFile)

    var install bool
    if err != nil && os.IsNotExist(err) {
        install = true
    }

    db, err := sql.Open("sqlite", dbFile)
    if err != nil {
        log.Fatalf("Ошибка подключения к базе данных: %v", err)
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
            log.Fatalf("Ошибка создания схемы базы данных: %v", err)
        }
        log.Println("База данных и таблица успешно созданы.")

        // Добавление тестовых данных
        insertQuery := `
        INSERT INTO scheduler (date, title, comment, repeat) 
        VALUES ('20250106', 'Тестовая задача 1', 'Описание задачи 1', 'd 7'),
               ('20250107', 'Тестовая задача 2', 'Описание задачи 2', 'd 14');
        `
        _, err = db.Exec(insertQuery)
        if err != nil {
            log.Fatalf("Ошибка вставки тестовых данных: %v", err)
        }
        log.Println("Тестовые данные добавлены в таблицу scheduler.")
    }
    return db
}