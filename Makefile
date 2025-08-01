build:
	docker compose build todo-app

run:
	docker compose up todo-app

test:
	go test -v ./...

migrate:
	migrate -path backend/migrations -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go