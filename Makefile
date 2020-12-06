test:
	docker-compose -f test-compose.yml up --d
	DEBUG=true go test

build:
	go build -o miauth cmd/main.go

run_integrated_test:
	./run_integrated_test.sh
