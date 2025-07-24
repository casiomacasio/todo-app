# 📝 TODO-приложение: REST API на Go + Frontend на React
Проект реализует полноценный backend на Go (с REST API, аутентификацией и middleware) и современный SPA frontend на React (Vite).

### Backend (Go):
- Разработка REST API с использованием `gin-gonic/gin`
- Архитектура по принципам Чистой Архитектуры
- Внедрение зависимостей и модульная структура
- Работа с базой данных PostgreSQL через `sqlx`
- Миграции базы данных через `migrate`
- Docker-среда разработки
- Конфигурация приложения через `spf13/viper` и `.env`
- JWT аутентификация и Refresh токены (хранятся в HttpOnly cookie)
- Middleware:
  - Проверка авторизации
  - Rate limiting:
    - Глобальный
    - На отдельного пользователя
- Swagger-документация через `swaggo/swag`
- Graceful shutdown

### Frontend (React + Vite):
- Одностраничное приложение (SPA)
- Маршрутизация через `react-router-dom`
- Авторизация с передачей cookie
- Интерактивный UI для управления задачами

## ⚙️ Как запустить проект

### 1. Клонировать репозиторий 
```bash
git clone https://github.com/casiomacasio/todo-app.git

cd todo-app
```
### 2. Создать .env файл и указать переменные окружения
Пример .env:
```bash
DB_PASSWORD=qwerty

REDIS_PASSWORD=redis

signingKey="jfklsdfj;eiwo;dskivewjieow;fiof"
```
### 3. Собрать и запустить backend
```bash
make build && make run
```
### 4. Применить миграции (при первом запуске)
```bash
make migrate
```
### 5. Запустить frontend
```bash
cd frontend

npm install

npm run dev
```
📚 Swagger-документация

После запуска backend вы можете открыть документацию API в браузере:
http://localhost:8000/swagger/index.html

📄 This documentation is also available in [English](README.en.md)

📄 Документация также доступна на [Английский](README.en.md)
