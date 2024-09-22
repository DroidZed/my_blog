#IMAGE_TAG=latest
include .env

run:
	go run cmd/my_blog/main.go

deps:
	go mod tidy

build:
	npx tailwindcss -i ./internal/asset/tailwind.css -o ./internal/asset/static/
	go-generate-fast
	go build -o ./bin/out.exe ./cmd/my_blog/

doc:
	swag init -g cmd/my_blog/main.go

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
	npx tailwindcss -i ./internal/asset/tailwind.css -o ./internal/asset/static/styles.css
	go-generate-fast

tools:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	go install github.com/oNaiPs/go-generate-fast@latest
