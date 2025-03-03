build:
	@go build -o bin/go-url-shortener

run: build
	@./bin/go-url-shortener

test:
	go test -v ./...
