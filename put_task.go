package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Обработчик для обновления задачи в базе данных
func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос выполнен методом PUT
	if r.Method != http.MethodPut {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	// Десериализуем JSON из тела запроса в структуру Task
	decoder := json.NewDecoder(r.Body)
	var task Task
	if err := decoder.Decode(&task); err != nil {
		// Если произошла ошибка десериализации, возвращаем ошибку 400
		http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, что указан id задачи
	if task.ID == "" {
		http.Error(w, `{"error":"не указан ID"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, что указан заголовок задачи
	if task.Title == "" {
		http.Error(w, `{"error":"не указан заголовок"}`, http.StatusBadRequest)
		return
	}

	// Получаем текущую дату в формате YYYYMMDD
	nowDate := time.Now().Format("20060102")
	// Если дата задачи не указана или указана как "today", устанавливаем текущую дату
	if task.Date == "" || task.Date == "today" {
		task.Date = nowDate
	} else {
		// Пробуем преобразовать строку с датой в формат time.Time
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"некорректный формат даты"}`, http.StatusBadRequest)
			return
		}
		// Если дата задачи раньше текущей, обрабатываем правило повторения
		if parsedDate.Before(time.Now()) {
			// Если повторение не указано, ставим дату на сегодня
			if task.Repeat == "" {
				task.Repeat = nowDate
			} else {
				// Если повторение указано, вычисляем следующую дату по правилу повторения
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					// Если правило повторения некорректно, возвращаем ошибку 400
					http.Error(w, `{"error":"некорректное правило повторения"}`, http.StatusBadRequest)
					return
				}
				// Если вычисленная дата не совпадает с текущей, обновляем её
				if task.Date != nowDate {
					task.Date = nextDate
				}
			}
		}
	}
	// Проверяем правило повторения
	if task.Repeat != "" {
		_, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}
	// Обновляем данные задачи в базе данных
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		// Если произошла ошибка при обновлении базы данных, возвращаем ошибку 500
		http.Error(w, `{"error": "ошибка обновления задачи в базе данных"}`, http.StatusInternalServerError)
		return
	}

	// Проверяем, было ли обновлено хотя бы одно поле
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error": "не удалось получить данные об обновленных полях"}`, http.StatusInternalServerError)
		return
	}
	// Если задача не найдена, возвращаем ошибку
	if rowsAffected == 0 {
		http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		return
	}

	// возвращаем успешный пустой JSON
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}