FROM golang:1.24-alpine

# Устанавливаем необходимые зависимости
RUN apk add --no-cache \
    bash \
    git \
    openssh \
    wireguard-tools \
    iproute2 \
    openresolv

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файл модулей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

EXPOSE 3000

# Запускаем main.go
CMD ["go", "run", "cmd/main/main.go"]
