# Используем официальный образ Go в качестве базового
FROM golang:1.22 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код в рабочую директорию
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Используем минимальный образ для запуска
FROM scratch

# Копируем собранный бинарный файл из стадии сборки
COPY --from=builder /app/main /main

# Указываем, что приложение будет слушать на порту 8080
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["./main"]