package main

import (
	"net/http"
	"encoding/json"
	"time"
	"strconv"
)

// Объявление структуры для ответа
type Response struct {
	Error string `json:"error,omitempty"`
	ID    string `json:"id,omitempty"`
}

// Обработчик для добавления новой задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос выполнен методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// Устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Декодируем JSON из тела запроса
	var task TaskCreate
	
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "ошибка десериализации"})
		return
	}
	// Проверка обязательных полей
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "не указан заголовок задачи"})
		return
	}
	// Получаем текущую дату
	nowDate := time.Now().Format("20060102")
	if task.Date == "" || task.Date == "today" {
		task.Date = nowDate
	} else {
		// Преобразуем дату в правильный формат
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Error: "некорректный формат даты"})
			return
		}
		// Используем функцию NextDate для вычисления следующей даты, если необходимо
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Date = nowDate
			} else {
				// Вычисляем следующую дату выполнения
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Response{Error: "ошибка при вычислении следующей даты"})
					return
				}

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
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}

	// Добавление задачи в базу данных
	db := createDatabase() // Используем уже существующую функцию для создания/получения базы данных
	defer db.Close()

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "ошибка добавления задачи в базу данных"})
		return
	}

	// Получаем id добавленной задачи
	id, err := res.LastInsertId()
	
	idResp := strconv.Itoa(int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "не удалось получить ID новой задачи"})
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": idResp})
}