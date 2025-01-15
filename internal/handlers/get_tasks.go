package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tandawg/agenda_project/internal/models"
)

// Константа для лимита количества задач в запросе
const TaskLimit = 50

// Обработчик для получения списка задач через маршрут /api/tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Проверяем, что запрос выполнен методом GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Считываем значение параметра "search" из строки запроса
	search := r.URL.Query().Get("search")
	var tasks []models.Task

	// Переменные для выполнения SQL-запроса
	var rows *sql.Rows
	var err error

	if search != "" {
		// Проверяем, является ли параметр "search" датой
		date, err := time.Parse("02.01.2006", search)
		if err == nil { // Если это дата
			// Преобразуем дату в формат базы данных (20060102)
			search = date.Format(models.DateFormat)
			// Выполняем запрос для поиска задач по дате
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
			rows, err = db.Query(query, search, TaskLimit)
			if err != nil {
				http.Error(w, `{"error": "Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
				return
			}
		} else { // Если параметр — строка для поиска в заголовке/комментарии
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
			rows, err = db.Query(query, "%"+search+"%", "%"+search+"%", TaskLimit)
			if err != nil {
				http.Error(w, `{"error": "Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
				return
			}
		}
	} else {
		// Если параметр "search" отсутствует, выбираем ближайшие задачи
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		rows, err = db.Query(query, TaskLimit)
		if err != nil {
			http.Error(w, `{"error": "Ошибка запроса к базе данных"}`, http.StatusInternalServerError)
			return
		}
	}

	// Обрабатываем строки результата запроса
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		// Считываем данные задачи из текущей строки результата
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error": "Ошибка чтения данных"}`, http.StatusInternalServerError)
			return
		}
		// Добавляем задачу в итоговый список
		tasks = append(tasks, task)
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error": "Ошибка при обработке строк результата"}`, http.StatusInternalServerError)
		return
	}

	// Если задачи не найдены, возвращаем пустой массив
	if len(tasks) == 0 {
		tasks = []models.Task{}
	}

	// Отправляем JSON-ответ с данными задач
	json.NewEncoder(w).Encode(map[string][]models.Task{"tasks": tasks})
}
