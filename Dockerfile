# Используем образ golang:alpine для сборки
FROM golang:latest

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все файлы из текущей директории в /app
COPY . .

# Соберём приложение
RUN go build -o /app/Bot ./main.go

# Запустим приложение
CMD ["/app/Bot"]