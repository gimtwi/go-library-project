build:
	@go build -o go-library-project

run: build
	@./go-library-project

test:
	@go test -v ./...