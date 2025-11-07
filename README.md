#  Delayed Notifier

**Delayed Notifier** — это микросервис для **отложенной отправки уведомлений** через Email и Telegram.  
Сервис поддерживает планирование уведомлений, проверку их статуса и отмену до момента отправки.  

```

##  Структура проекта


delayed_notifier/
├── cmd/
│ └── main.go # Точка входа в приложение
├── db/
│ └── dumps/ # Директория для дампов (резервные данные)
├── internal/
│ ├── app/
│ │ └── app.go # Инициализация, запуск и конфигурация приложения
│ ├── cache/
│ │ └── cache.go # Работа с Redis
│ ├── handlers/
│ │ ├── addNotice.go # POST /notify — создание уведомления
│ │ ├── deleteNotice.go # DELETE /notify/:id — удаление уведомления
│ │ ├── getStatus.go # GET /notify/:id — получение статуса уведомления
│ │ └── router.go # Регистрация всех маршрутов (включая web UI)
│ ├── model/
│ │ └── model.go # Определение структур уведомлений и статусов
│ ├── notifier/
│ │ ├── email/
│ │ │ └── email.go # Отправка уведомлений по Email
│ │ └── telegram/
│ │ └── telegram.go # Отправка уведомлений в Telegram
│ ├── rabbitmq/
│ │ └── rabbit.go # Работа с RabbitMQ (очередь и задержка)
│ ├── service/
│ │ ├── data_deleter/
│ │ │ └── service.go # Удаление уведомлений
│ │ ├── data_saver/
│ │ │ └── service.go # Сохранение уведомлений
│ │ └── status_getter/
│ │ └── service.go # Получение статуса уведомления
│ └── storage/
│ └── storage.go # Работа с хранилищем (интеграция Redis + кеш)
├── temp/ # Временные файлы, логи
├── web/
│ └── index.html # Простой UI для тестирования (HTML + JS)
├── Dockerfile
├── docker-compose.yml
└── README.md


```

##  Основные возможности

✅ Планирование уведомлений с точным временем (`send_at`)  
✅ Поддержка отправки через **Email** и **Telegram**  
✅ Хранение данных и статусов в **Redis**  
✅ Очередь задач и отложенная отправка через **RabbitMQ**  
✅ Повторные попытки отправки с экспоненциальной задержкой  
✅ Простой **веб-интерфейс** (`web/index.html`) для тестирования  

---

##  Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/Vladimirmoscow84/delayed_notifier.git
cd delayed_notifier

2. Создайте .env в корне проекта
SERVER_ADDRESS=:7540
REDIS_URI=redis://redis:6379
RABBIT_URI=amqp://user:password@rabbitmq:5672/
EMAIL_HOST=smtp.example.com
EMAIL_PORT=587
EMAIL_USER=user@example.com
EMAIL_PASS=password
EMAIL_FROM=user@example.com
EMAIL_TO=recipient@example.com
TELEGRAM_TOKEN=123456:ABC-DEF
TELEGRAM_CHAT_ID=987654321

Укажите реальные SMTP-данные, чтобы тестировать отправку писем.
Telegram можно отключить, оставив токен и chat_id пустыми.

3. Запуск через Docker Compose
docker-compose build --no-cache
docker-compose up -d

 После запуска:
| Компонент   | Адрес                                                        | Назначение                      |
| ----------- | ------------------------------------------------------------ | ------------------------------- |
| Web UI      | [http://localhost:7540](http://localhost:7540)               | Создание и просмотр уведомлений |
| RabbitMQ UI | [http://localhost:15672](http://localhost:15672)             | Очереди и сообщения             |
| Redis       | localhost:6379                                               | Хранилище уведомлений           |
| API         | [http://localhost:7540/notify](http://localhost:7540/notify) | REST API                        |

Веб-интерфейс

Путь: web/index.html
Открыть в браузере: http://localhost:7540

Функции интерфейса:

Создать уведомление (POST /notify)

Проверить статус (GET /notify/:id)

Удалить уведомление (DELETE /notify/:id)

 REST API
▶POST /notify — создание уведомления

Пример запроса:
{
  "subject": "Напоминание",
  "body": "Позвонить клиенту в 15:00",
  "send_at": "2025-10-21T15:00:00Z"
}
Ответ:
{
  "id": "1",
  "status": "scheduled",
  "send_at": "2025-10-21T15:00:00Z"
}
 GET /notify/{id} — статус уведомления

Ответ:
{
  "id": "1",
  "status": "scheduled",
  "subject": "Напоминание",
  "body": "Позвонить клиенту в 15:00",
  "send_at": "2025-10-21T15:00:00Z"
}
 DELETE /notify/{id} — отмена уведомления

Ответ:
{
  "id": "1",
  "status": "deleted"
}

Принцип работы
1️⃣ Клиент (UI / API) отправляет POST /notify с временем отправки (send_at)
2️⃣ Сервис сохраняет уведомление в Redis
3️⃣ Уведомление публикуется в RabbitMQ с отложенной доставкой
4️⃣ Когда наступает время, воркер извлекает сообщение
5️⃣ Сервис отправляет уведомление по Email или Telegram
6️⃣ Статус обновляется в Redis

 Технологии
| Компонент            | Назначение                       |
| -------------------- | -------------------------------- |
| **Go (Golang)**      | Основной язык                    |
| **Redis**            | Хранилище уведомлений и статусов |
| **RabbitMQ**         | Очередь задач с TTL (задержкой)  |
| **SMTP (Email)**     | Отправка писем                   |
| **Telegram Bot API** | Отправка сообщений в Telegram    |
| **Gin**              | HTTP API                         |
| **Docker Compose**   | Контейнеризация всех сервисов    |
| **HTML + JS**        | Веб-интерфейс                    |

Примеры запросов через curl
Создание уведомления
curl -X POST http://localhost:7540/notify \
-H "Content-Type: application/json" \
-d '{"subject":"Тест","body":"Сообщение","send_at":"2025-10-21T15:00:00Z"}'

Проверка статуса
curl http://localhost:7540/notify/1

Удаление уведомления
curl -X DELETE http://localhost:7540/notify/1

Дополнительно

Логи можно писать в /temp/

Все временные файлы — в /temp/

Redis можно сбросить командой:
docker exec -it redis_server redis-cli flushall

Архитектура проекта
                    ┌────────────────────┐
                    │   Web UI / API     │
                    │  (index.html, /notify) │
                    └──────────┬─────────┘
                               │ POST /notify
                               ▼
                      ┌──────────────────┐
                      │   Handlers       │
                      │  (addNotice etc) │
                      └──────────┬────────┘
                                 │
                                 ▼
                        ┌───────────────┐
                        │    Redis      │
                        │  (Storage)    │
                        └──────┬────────┘
                               │
                               ▼
                       ┌─────────────┐
                       │  RabbitMQ   │
                       │  (Queue + TTL) │
                       └──────┬────────┘
                              │ (по истечении TTL)
                              ▼
                ┌────────────────────────────┐
                │   Notifier Service         │
                │  (Email + Telegram Sender) │
                └──────────┬─────────────────┘
                           │
                           ▼
                     Получатель уведомления



Лицензия
MIT © 2025 — VladimirMoscow84
