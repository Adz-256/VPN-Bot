# VPN Бот (Telegram + WireGuard)

Telegram-бот для продажи доступа к VPN через WireGuard. После оплаты через ЮMoney бот автоматически выдаёт конфигурационный файл и QR-код для подключения.

## 🛠 Возможности

- Покупка VPN внутри Telegram
- Поддержка WireGuard (удобен для мобильных и десктопных клиентов)
- Автоматическая выдача:
  - `.conf` файла
  - QR-кода
- Актуальные тарифы: 1 / 3 / 6 месяцев
- Оплата через **ЮMoney**
- Хранение данных в PostgreSQL
- Webhook-подтверждение оплаты через [smee.io](https://smee.io)
- Сервер: 🇳🇱 Нидерланды

## ⚙️ Стек технологий

- Язык: Go
- Telegram SDK: [`github.com/go-telegram/bot`](https://github.com/go-telegram/bot)
- VPN: [WireGuard](https://www.wireguard.com/)
- БД: PostgreSQL
- Webhook-прокси: [smee.io](https://smee.io)

## 🚀 Запуск

1. Создайте файл `.env` по следующему примеру:

    ```env
    DSN=postgresql://postgres:postgres@postgres:5432/cheapvpn
    ENV=development

    BOT_TOKEN=8049404870:AAF-4MRbvwEnh4oTMafto7me8bCE-yWpowI
    PAYMENT_ACCOUNT=4100117034899495

    WIREGUARD_CONFIG_PATH=wg0.conf
    WIREGUARD_ADDRESS=185.244.48.22
    WIREGUARD_INTERFACE_NAME=wg0
    WIREGUARD_OUT=config
    WIREGUARD_PORT=51820

    WEBHOOK_ADDRESS=127.0.0.1
    WEBHOOK_PORT=3000
    ```

2. Убедитесь, что WireGuard установлен и настроен.

3. Запустите smee-прокси (в отдельном терминале/сессии, вне Docker):

    ```bash
    npx smee -u https://smee.io/your-channel --target http://localhost:3000/webhook
    ```

4. Запустите PostgreSQL, если используется Docker и примените миграции:

    ```bash
    docker-compose up -d postgres
    ```

5. Соберите и запустите бота:

    ```bash
    go mod tidy
    go run main.go
    ```

  ```bash
    goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up
  ```

## 💳 Тарифы

Задаются в базе. Есть предустановленные в файлах миграций

## 🔐 Подключение через WireGuard

Пользователь получает:
- Файл конфигурации `.conf`
- QR-код для быстрого подключения

WireGuard не требует логина и пароля. Для подключения просто импортируйте файл или отсканируйте QR-код в [официальном приложении](https://www.wireguard.com/install/).

## 🏗 Планы на будущее

- Серверы в других странах
- Поддержка криптооплаты
- Панель управления VPN-подписками
- Dockerized smee-прокси
