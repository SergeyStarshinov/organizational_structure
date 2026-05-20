export CONFIG_PATH=./config/local.yaml

.PHONY: run
run: 
	go run ./cmd/main.go

.PHONY: test
test:
	go test ./test
