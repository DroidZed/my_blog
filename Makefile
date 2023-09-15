#IMAGE_TAG=latest
include .env

init_deps:
	rm go.mod go.sum go.work.sum
	go mod init github.com/DroidZed/go_lance
	go mod tidy

build:
	go build -v -o bin/golance cmd/go_lance/main.go

dev:
	go run cmd/go_lance/main.go

doc:
	swag init -g cmd/go_lance/main.go

prod:
	./bin/golance.exe

dockerImage:
	docker build -t droidzed/golance:$(IMAGE_TAG) .

compose:
	docker compose up -d

dockerUp: dockerImage compose