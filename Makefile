.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans

down:
	docker-compose down -v

lint:
	golangci-lint run

swag:
	swag init -g cmd/app/main.go
