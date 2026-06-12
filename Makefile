# Application names and configuration
APP_NAME=edns
DB_URL=postgres://admin:admin@localhost:5432/notifications?sslmode=disable

run-api:
	go run api-service/cmd/server/main.go

build-api:
	go build -o bin/$(APP_NAME) api-service/cmd/server/main.go

build-worker:
	go build -o bin/$(APP_NAME)-worker worker-service/cmd/worker/main.go

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

docker-up-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

docker-down-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

docker-up-prod:
	docker compose -f docker-compose.yml up --build

docker-down-prod:
	docker compose -f docker-compose.yml down

# database migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations $(name)

# manage dependencies across all services
deps:
	go mod tidy

clean:
	rm -rf bin
