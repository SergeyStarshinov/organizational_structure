export CONFIG_PATH=./config/local.yaml

.PHONY: run
run: 
	docker pull postgres:latest
	docker compose up -d

.PHONY: test
test:
	go test ./internal/web
