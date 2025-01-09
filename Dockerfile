# Этап сборки
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod download

# Собираем исполняемый файл для Linux
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

# Финальный образ
FROM ubuntu:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный исполняемый файл из первого этапа
COPY --from=builder /app/main /app/

# Копируем директорию web из первого этапа
COPY --from=builder /app/agenda_project/web /app/agenda_project/web

# Настраиваем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

# Указываем порт веб-сервера
EXPOSE 7540

# Команда для запуска веб-сервера
CMD ["/app/main"]