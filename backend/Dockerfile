FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-app ./backend/cmd/main.go

CMD ["./todo-app"]