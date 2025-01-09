package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate рассчитывает следующую дату задачи с учетом правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Преобразуем строку с датой задачи в формат time.Time
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %s", date)
	}

	// Проверяем, что год даты задачи находится в допустимом диапазоне (1900-2100)
	if taskDate.Year() < 1900 || taskDate.Year() > 2100 {
		fmt.Printf("Отказ: год %d вне допустимого диапазона (1900–2100)\n", taskDate.Year())
		return "", nil // Возвращаем пустую строку для некорректных годов
	}

	// Разбиение правила повторения на части
	repeatParts := strings.Fields(repeat)
	if len(repeatParts) == 0 {
		// Если правило пустое, возвращаем пустую строку
		return "", nil
	}

	// Обрабатываем правило повторения
	switch repeatParts[0] {
	case "d": // Повторение в днях
		if len(repeatParts) != 2 {
			// Если параметров меньше двух (не указано количество дней), возвращаем пустую строку
			return "", nil
		}
		// Преобразуем количество дней в целое число
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days <= 0 || days > 400 {
			fmt.Printf("Отказ: некорректное количество дней %s\n", repeatParts[1])
			return "", nil
		}
		// Вычисляем следующую дату, добавляя дни
		nextDate := taskDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			// Если дата ещё раньше текущей, добавляем дни снова
			nextDate = nextDate.AddDate(0, 0, days)
			// Если дата выходит за пределы года 2100, прекращаем
			if nextDate.Year() > 2100 {
				return "", nil
			}
		}
		// Возвращаем следующую дату в формате YYYYMMDD
		return nextDate.Format("20060102"), nil

	case "y": // Повторение в годах
		// Добавляем один год к начальной дате задачи
		nextDate := taskDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			// Если дата раньше текущей, добавляем год снова
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		if nextDate.Year() > 2100 {
			return "", nil
		}
		return nextDate.Format("20060102"), nil

	default:
		// Если правило повторения не поддерживается, возвращаем пустую строку
		return "", nil
	}
}

// Вспомогательная функция, проверяет корректность формата даты в строке (YYYYMMDD)
func isValidDate(date string) bool {
	// Проверяем, что длина строки равна 8 символам
	if len(date) != 8 {
		return false
	}
	// Извлекаем год, месяц и день из строки
	year := date[:4]
	month := date[4:6]
	day := date[6:8]

	// Преобразуем год, месяц и день в целые числа
	yearNum, err1 := strconv.Atoi(year)
	monthNum, err2 := strconv.Atoi(month)
	dayNum, err3 := strconv.Atoi(day)

	// Если преобразование не удалось, возвращаем false
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// Проверка по верхней границе дат, включая исторические даты
	if yearNum > 2100 || monthNum < 1 || monthNum > 12 || dayNum < 1 || dayNum > 31 {
		return false
	}

	return true
}