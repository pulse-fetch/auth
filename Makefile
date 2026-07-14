build:
	go build -o bin/auth cmd/main.go
run: build

	./bin/auth --config=./configs/local.yaml
run-postgres:
	docker compose down -v
	docker compose up -d