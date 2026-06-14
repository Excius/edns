FROM golang:1.26.1-alpine3.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ws-worker ./websocket-service/cmd/worker

EXPOSE 8081

CMD ["./worker"]
