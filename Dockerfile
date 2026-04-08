FROM golang:1.25-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

RUN go build -o todo-api ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-api .
COPY config/ config/

EXPOSE 8080

CMD ["./todo-api"]