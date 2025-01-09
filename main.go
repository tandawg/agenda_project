package main

func main() {
	// Запуск базы данных
	DB := createDatabase()
	defer DB.Close()

	// Запуск веб-сервера
	startServer()
}