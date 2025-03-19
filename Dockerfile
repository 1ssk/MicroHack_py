# Используем официальный образ Go 1.22.5 как базовый
FROM golang:1.22.5-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходники проекта в контейнер
COPY . .

# Собираем приложение с оптимизацией размера бинарного файла
RUN go build -ldflags="-s -w" -o main ./cmd/main.go

# Используем минимальный образ для продакшн-версии
FROM alpine:latest  

# Устанавливаем рабочую директорию в финальном контейнере
WORKDIR /app

# Копируем собранное приложение из builder-контейнера
COPY --from=builder /app/main /app/main

# Копируем статические файлы и шаблоны
COPY ./web/static /app/web/static
COPY ./web/templates /app/web/templates

# Копируем .env файл (если он есть) в контейнер
COPY .env .env

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"]