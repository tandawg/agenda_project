package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tandawg/agenda_project/internal/models"
	"github.com/tandawg/agenda_project/internal/services"
)

// Функция обработчика для маршрута /api/nextdate
// Вычисляет следующую дату на основе заданного правила повторения
func NextDateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Извлекаем параметры запроса: текущую дату (now), начальную дату (date) и правило повторения (repeat)
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Преобразуем параметр now в формат time.Time
	now, err := time.Parse(models.DateFormat, nowStr)
	if err != nil {
		// Если формат даты некорректен, возвращаем ошибку
		http.Error(w, "Некорректный формат даты 'now'", http.StatusBadRequest)
		return
	}

	// Пробуем преобразовать параметр date в формат time.Time
	dateParsed, err := time.Parse(models.DateFormat, date)
	if err != nil || !services.IsValidDate(date) { // Используем isValidDate
		// Если date некорректна или невалидна, возвращаем пустую строку
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "")
		return
	}

	// Обрабатываем правило повторения repeat
	switch {
	case repeat == "y":
		// Добавляем годы, пока не достигнем ближайшей даты после now
		nextDate := dateParsed.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		// Отправляем результат в формате YYYYMMDD
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, nextDate.Format(models.DateFormat))

	case strings.HasPrefix(repeat, "d "):
		// Если правило указывает повторение через определённое количество дней
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days <= 0 || days > 400 {
			// // Если количество дней некорректно, возвращаем пустую строку
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, "")
			return
		}
		// Рассчитываем ближайшую дату, добавляя дни к начальной дате
		nextDate := dateParsed.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, nextDate.Format(models.DateFormat))

	default:
		// Если правило повторения отсутствует или содержит некорректное значение
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "")
	}
}
