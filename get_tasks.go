package main

import (
	"encoding/json"
	"net/http"
	"time"
	"database/sql"
)

// Обработчик для получения списка задач через маршрут /api/tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос выполнен методом GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Считываем значение параметра "search" из строки запроса
search := r.URL.Query().Get("search")
var tasks []Task

// Переменные для выполнения SQL-запроса
var rows *sql.Rows
var err error

if search != "" {
    // Проверяем, является ли параметр "search" датой
    date, err := time.Parse("02.01.2006", search)
    if err == nil { // Если это дата
        // Преобразуем дату в формат базы данных (20060102)
        search = date.Format("20060102")
        // Выполняем запрос для поиска задач по дате
        rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT 50`, search)
        if err != nil {
            http.Error(w, `{"error": "Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
            return
        }        
    } else { // Если параметр — строка для поиска в заголовке/комментарии
        rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT 50`, "%"+search+"%", "%"+search+"%")
        if err != nil {
            http.Error(w, `{"error": "Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
            return
        }        
    }
} else {
    // Если параметр "search" отсутствует, выбираем ближайшие задачи
    rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50`)
    if err != nil {
        http.Error(w, `{"error": "Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
        return
    }    
}

// Обрабатываем строки результата запроса
defer rows.Close()
for rows.Next() {
    var task Task
    // Считываем данные задачи из текущей строки результата
    if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка чтения данных"})
        return
    }
    // Добавляем задачу в итоговый список
    tasks = append(tasks, task)
}

// Если задачи не найдены, возвращаем пустой массив
if len(tasks) == 0 {
    tasks = []Task{}
}

// Отправляем JSON-ответ с данными задач
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string][]Task{"tasks": tasks})
}