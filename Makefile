.PHONY: build run test

# Build the Go application
build:
	go build -o main handler/handler.go

# Run the Go application
run:
	go run main.go

# Test the Go application
test:
	go test ./...

dev:
	vercel

deploy:
	vercel --prod

doc:
	swag init

run-doc:
	node server.js