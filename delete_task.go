package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Обработчик для удаления задачи из базы данных
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос выполнен методом DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	// Получаем идентификатор задачи из пути/строки запроса
	idTask := chi.URLParam(r, "id")
	if idTask == "" {
		idTask = r.URL.Query().Get("id")
	}

	// Если id отсутствует, возвращаем ошибку
	if idTask == "" {
		http.Error(w, `{"error": "не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Подготавливаем запрос на удаление задачи по её id
	deleteQuery := "DELETE FROM scheduler WHERE id = ?"
	res, err := DB.Exec(deleteQuery, idTask)
	if err != nil {
		http.Error(w, `{"error": "ошибка удаления задачи"}`, http.StatusInternalServerError)
		return
	}

	// Проверяем, затронул ли запрос удаления хотя бы одну строку
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error": "ошибка проверки удаления задачи"}`, http.StatusInternalServerError)
		return
	}

	// Если удаление не затронуло строки, задача с таким идентификатором не найдена
	if rowsAffected == 0 {
		http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Если задача успешно удалена, отправляем пустой успешный ответ
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}