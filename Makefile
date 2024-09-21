#IMAGE_TAG=latest
include .env

run:
	go run cmd/my_blog/main.go

deps:
	go mod tidy

build:
	docker build -t droidzed/blog:$(IMAGE_TAG) .

doc:
	swag init -g cmd/my_blog/main.go

compose:
	docker compose up -d

decompose:
	docker compose down

dev:
	air

prod:
	./bin/golance

module:
	mkdir ./internal/${DIR}
	echo "package ${DIR}" > ./internal/${DIR}/controller.go
	echo "package ${DIR}" > ./internal/${DIR}/service.go
	echo "package ${DIR}" > ./internal/${DIR}/models.go

templates:
	templ generate -lazy

tools:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
