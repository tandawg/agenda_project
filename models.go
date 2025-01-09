package main

// TaskCreate используется для создания новой задачи
type TaskCreate struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Task используется для получения данных из базы
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}