package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Обработчик для получения информации о задаче по её id
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос выполнен методом GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	// Получаем id задачи из пути/строки запроса
	idTask := chi.URLParam(r, "id")
	if idTask == "" {
		idTask = r.URL.Query().Get("id")
	}

	// Если id задачи не указан, возвращаем сообщение об ошибке
	if idTask == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// SQL-запрос для поиска задачи в базе данных
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := DB.QueryRow(query, idTask)

	// Переменная для хранения данных найденной задачи
	var task Task
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Если задача не найдена, возвращаем соответствующий код ошибки
			http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		} else {
			// При другой ошибке возвращаем сообщение о сбое работы с базой
			http.Error(w, `{"error": "ошибка получения задачи из базы данных"}`, http.StatusInternalServerError)
		}
		return
	}

	// Формируем и отправляем JSON-ответ с информацией о задаче
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		// Если возникла ошибка формирования JSON, возвращаем соответствующий код
		http.Error(w, `{"error":"ошибка формирования JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
}