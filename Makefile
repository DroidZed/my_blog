#IMAGE_TAG=latest
include .env

deps:
	go mod tidy

build:
	docker build -t droidzed/golance:$(IMAGE_TAG) .

watch:
	air

doc:
	swag init -g cmd/go_lance/main.go

compose:
	docker compose up -d

decompose:
	docker compose down

dev:
	build compose

prod:
	./bin/golance.exe

module:
	mkdir ./internal/${DIR}
	echo "package ${DIR}" > ./internal/${DIR}/controller.go
	echo "package ${DIR}" > ./internal/${DIR}/service.go
	echo "package ${DIR}" > ./internal/${DIR}/models.go
