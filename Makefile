test:
	docker-compose -f test-compose.yml up --d
	DEBUG=true go test
build:
	go build -o miauth cmd/main.go