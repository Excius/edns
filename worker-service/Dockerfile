FROM golang:1.25.8-alpine3.23

WORKDIR /app

COPY ./go.mod ./
RUN go mod download

COPY . .

RUN go build -o worker ./cmd/worker

CMD ["./worker"]