#IMAGE_TAG=latest
include .env

run:
	go run cmd/api/main.go

deps:
	go mod tidy

build:
	go build -o ./bin/out.exe ./cmd/api/

doc:
	swag init -g cmd/api/main.go

dev:
	air

clean:
	rm ./bin/out.exe

module:
	mkdir ./internal/${DIR}
	echo "package ${DIR}" > ./internal/${DIR}/controller.go
	echo "package ${DIR}" > ./internal/${DIR}/service.go
	echo "package ${DIR}" > ./internal/${DIR}/models.go

templates:
	npx tailwindcss -i ./cmd/web/assets/css/input.css -o ./cmd/web/assets/css/output.css
	go-generate-fast

tools:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	go install github.com/oNaiPs/go-generate-fast@latest
