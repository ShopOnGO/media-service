FROM golang:1.23.3 AS builder

WORKDIR /media

# Отключаем CGO для статической компиляции
 ENV CGO_ENABLED=0

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download && go mod verify

# Копируем весь код
COPY . .

# Компилируем бинарник
RUN go build -o /media/media_service ./cmd/server.go



# Второй этап: финальный образ (без лишних инструментов)
FROM alpine:latest

WORKDIR /media

COPY .env /media/.env

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /media/media_service /media/media_service

# Запуск приложения
CMD ["/media/media_service"]
