run_api:
	go run cmd/go_lance/main.go

init_deps:
	rm -f go.mod go.sum
	go mod init github.com/DroidZed/go_lance
	go mod tidy