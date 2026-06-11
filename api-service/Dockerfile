FROM golang:1.26.1-alpine3.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./api-service/cmd/server

EXPOSE 8080

CMD ["./api"]