APP_NAME=notification-api

run-api:
	cd api-service && go run cmd/server/main.go

build-api:
	cd api-service && go build -o bin/$(APP_NAME) cmd/server/main.go

test:
	cd api-service && go test ./...

lint:
	cd api-service && golangci-lint run

fmt:
	cd api-service && go fmt ./...

vet:
	cd api-service && go vet ./...

docker-up-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

docker-down-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

deps:
	cd api-service && do mod tidy

clean:
	cd api-service && rm -rf bin
