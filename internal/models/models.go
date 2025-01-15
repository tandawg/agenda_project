package models

// Константа для формата даты
const DateFormat = "20060102"

// Task используется для создания новой задачи, а также для получения данных из базы
type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
