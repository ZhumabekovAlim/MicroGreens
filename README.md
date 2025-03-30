
# 🌱 MicroGreens

**MicroGreens** — backend-система на Go для агропроектов, предназначенная для управления выращиванием микрозелени. Поддерживает партии, ежедневные наблюдения, аналитику, уведомления и AI-интеграцию.

---

## 🚀 Возможности

### 🔐 Аутентификация
- Регистрация с безопасным хешированием пароля (`bcrypt`)
- Логин пользователя с возвратом `user_id`

### 🌿 Управление посевами
- CRUD микрозелени `/api/microgreens`
- CRUD партий выращивания `/api/batches`

### 📆 Ежедневные наблюдения
- Ведение записей `/api/observations`
- Загрузка фотоотчётов `/api/photos`

### 📊 Аналитика
- Прогресс посева (оставшиеся дни и %)
- Влажность и рост за последние 7 дней
- Анализ от AI (Success Rate, Est. Yield, Quality)
- Наблюдения, пропущенные сегодня

### 🤖 AI-чат
- WebSocket чат с AI `/ws/ai` (ChatGPT API)

### 🔔 Уведомления (FCM)
- Добавление токенов устройств
- Отправка напоминаний вручную и по расписанию
- История отправленных уведомлений

---

## 🛠️ Технологии

- Go 1.20+
- MySQL / MariaDB
- Firebase Cloud Messaging
- OpenAI GPT API
- WebSocket (gorilla)
- bcrypt для авторизации

---

## 📂 Архитектура проекта

```
/cmd/web              # Запуск сервера, инициализация
/internal/handlers    # Обработчики HTTP и WS
/internal/services    # Бизнес-логика
/internal/repositories # Работа с БД
/internal/models      # Структуры данных
```

---

## ▶️ Быстрый старт

```bash
git clone https://github.com/yourname/microgreens.git
cd microgreens

cp .env.example .env        # Укажи ключи OpenAI и Firebase
go run ./cmd/web            # Запуск сервера
```

---

## 📩 Пример логина

```http
POST /api/login
{
  "email": "user@example.com",
  "password": "123456"
}
```

**Ответ:**
```json
{ "user_id": 1 }
```



