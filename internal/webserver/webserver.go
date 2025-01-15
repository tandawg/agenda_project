package webserver

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/tandawg/agenda_project/internal/handlers"
)

// startServer запускает HTTP-сервер для обработки запросов API
func StartServer(db *sql.DB) {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	staticPath := "./web"
	fmt.Println("Serving static files from:", staticPath)

	// Обработчик для маршрута корневого каталога (статические файлы)
	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	// Обработчики API с передачей объекта db
	http.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		handlers.NextDateHandler(w, r, db)
	})
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTaskHandler(w, r, db)
	})
	http.HandleFunc("/api/gettask", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTaskHandler(w, r, db)
	})
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTasksHandler(w, r, db)
	})
	http.HandleFunc("/api/puttask", func(w http.ResponseWriter, r *http.Request) {
		handlers.PutTaskHandler(w, r, db)
	})
	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		handlers.DoneTaskHandler(w, r, db)
	})
	http.HandleFunc("/api/deletetask", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTaskHandler(w, r, db)
	})

	fmt.Printf("Сервер работает на порту %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
