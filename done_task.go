package main

import (
	"database/sql"
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
)

// Обработчик для отметки задачи как выполненной
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос выполнен методом POST
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	// Получаем id задачи из пути/ строки запроса
	idTask := chi.URLParam(r, "id")
	if idTask == "" {
		idTask = r.URL.Query().Get("id")
	}

	// Если id не указан, возвращаем ошибку
	if idTask == "" {
		http.Error(w, `{"error": "не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Ищем задачу в базе данных по её id
	var task Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := DB.QueryRow(query, idTask)
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		// Если задача не найдена, возвращаем соответствующую ошибку
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		// Если возникла ошибка доступа к базе данных
		} else {
			http.Error(w, `{"error": "ошибка получения задачи из БД"}`, http.StatusInternalServerError)
		}
		return
	}

	// Проверяем, является ли задача одноразовой или повторяющейся
	if task.Repeat == "" {
		// Для одноразовой задачи выполняем её удаление из базы данных
		deleteQuery := "DELETE FROM scheduler WHERE id = ?"
		_, err := DB.Exec(deleteQuery, idTask)
		if err != nil {
			http.Error(w, `{"error": "ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// Для повторяющейся задачи рассчитываем дату следующего выполнения
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error": "ошибка расчёта следующей даты"}`, http.StatusInternalServerError)
			return
		}

		// Обновляем дату выполнения задачи в базе данных
		updateQuery := "UPDATE scheduler SET date = ? WHERE id = ?"
		_, err = DB.Exec(updateQuery, nextDate, idTask)
		if err != nil {
			http.Error(w, `{"error": "ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	// Если выполнение успешно, возвращаем пустой успешный JSON
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}