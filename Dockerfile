FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV TODO_PORT=7540

ENV TODO_DBFILE=/app/scheduler.db

EXPOSE ${TODO_PORT}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

CMD ["/my_app"] 