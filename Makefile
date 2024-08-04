#IMAGE_TAG=latest
include .env

run:
	go run cmd/my_blog/main.go -o bin/my_blog

deps:
	go mod tidy

build:
	docker build -t droidzed/golance:$(IMAGE_TAG) .

doc:
	swag init -g cmd/my_blog/main.go

compose:
	docker compose up -d

decompose:
	docker compose down

dev:
	build compose

prod:
	./bin/golance

module:
	mkdir ./internal/${DIR}
	echo "package ${DIR}" > ./internal/${DIR}/controller.go
	echo "package ${DIR}" > ./internal/${DIR}/service.go
	echo "package ${DIR}" > ./internal/${DIR}/models.go

templates:
	templ generate

setup-dev-container:
	sudo apt-get update && sudo apt-get upgrade -y
	curl -fsSL https://bun.sh/install | bash
	go install github.com/a-h/templ/cmd/templ@latest
