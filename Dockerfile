FROM golang:1.22.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o todo-app

FROM alpine:latest

# Установка переменных окружения по умолчанию
ENV TODO_PORT=7540 \
    TODO_DBFILE=/app/scheduler.db \
    TODO_PASSWORD=

WORKDIR /app

COPY --from=build /app/todo-app ./
COPY --from=build /app/web ./web
COPY --from=build /app/.env ./

EXPOSE 7540

CMD ["./todo-app"]
