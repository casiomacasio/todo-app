# 📝 TODO App: REST API in Go + Frontend in React
This project implements a full-featured backend in Go (with REST API, authentication, and middleware) and a modern SPA frontend in React (Vite).

### Backend (Go):
- Development of a REST API using `gin-gonic/gin`
- Architecture based on Clean Architecture principles
- Dependency injection and modular structure
- Work with PostgreSQL database via `sqlx`
- Database migrations using `migrate`
- Docker development environment
- Application configuration via `spf13/viper` and `.env`
- JWT authentication and refresh tokens (stored in HttpOnly cookies)
- Middleware:
  - Authorization check
  - Rate limiting:
    - Global
    - Per user
- Swagger documentation via `swaggo/swag`
- Graceful shutdown

### Frontend (React + Vite):
- Single Page Application (SPA)
- Routing via `react-router-dom`
- Authorization using cookies
- Interactive UI for task management

## ⚙️ How to Run the Project

### 1. Clone the repository
```bash
git clone https://github.com/casiomacasio/todo-app.git
cd todo-app
```

### 2. Create a `.env` file and set environment variables
```bash
Example .env:
DB_PASSWORD=qwerty
REDIS_PASSWORD=redis
signingKey="jfklsdfj;eiwo;dskivewjieow;fiof"
```
### 3. Build and run the backend
```bash
make build && make run
```
### 4. Apply migrations (on first run)
```bash
make migrate
```
### 5. Run the frontend
```bash
cd frontend
npm install
npm run dev
```
📚 **Swagger Documentation**
After starting the backend, you can open the API documentation in your browser:
http://localhost:8000/swagger/index.html  
📄 This documentation is also available in [Russian](README.md)   

📄 Документация также доступна на [Русский](README.en.md)
