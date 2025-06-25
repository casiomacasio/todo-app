FROM golang:latest

RUN go version
COPY ./ ./

# RUN apt-get update
# RUN apt-get -y install postgresql-client

# RUN chmod +x wait-for-postgres.sh

RUN go mod download

RUN go build -o todo-app ./cmd/main.go

CMD ["./todo-app"]