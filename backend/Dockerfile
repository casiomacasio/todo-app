FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-app ./backend/cmd/main.go

RUN mkdir -p /app/configs && cp ./backend/configs/config.yaml /app/configs/

CMD ["./todo-app"]