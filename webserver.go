package main

import (
	"fmt"
	"net/http"
	"os"
)

// startServer запускает HTTP-сервер для обработки запросов API
func startServer() {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	staticPath := "./web"
	fmt.Println("Serving static files from:", staticPath)

	// Обработчик для маршрута корневого каталога (статические файлы)
	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	// Обработчик для вычисления следующей даты
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Обработчик для добавления новой задачи (POST)
	http.HandleFunc("/api/task", AddTaskHandler)

	// Обработчик для получения данных задачи (GET)
	http.HandleFunc("/api/gettask", GetTaskHandler)

	// Обработчик для получения списка задач (GET)
	http.HandleFunc("/api/tasks", GetTasksHandler)

	// Обработчик для обновления данных задачи (PUT)
	http.HandleFunc("/api/puttask", PutTaskHandler)

	// Обработчик для завершения задачи (POST)
	http.HandleFunc("/api/task/done", DoneTaskHandler)

	// Обработчик для удаления задачи (DELETE)
	http.HandleFunc("/api/deletetask", DeleteTaskHandler)

	fmt.Printf("Сервер работает на порту %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}