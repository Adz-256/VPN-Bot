FROM golang:1.24-alpine

# Устанавливаем необходимые зависимости
RUN apk add --no-cache \
    bash \
    git \
    openssh \
    wireguard-tools \
    iproute2 \
    iptables \
    openresolv \
    iputils

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файл модулей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

USER root
# Выполняем команды iptables для настройки сети
ENTRYPOINT ["sh", "-c", "iptables -t nat -A POSTROUTING -s 10.9.0.0/32 -o eth0 -j MASQUERADE && \
    iptables -A INPUT -p udp -m udp --dport 51820 -j ACCEPT && \
    iptables -A FORWARD -i wg1 -j ACCEPT && \
    iptables -A FORWARD -o wg1 -j ACCEPT && \
    iptables -t nat -A PREROUTING -p udp --dport 51825 -j REDIRECT --to-port 51820 && \
    exec go run cmd/main/main.go"]

