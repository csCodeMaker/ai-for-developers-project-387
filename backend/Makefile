.PHONY: gen run build test tidy

# Генерация типов из OpenAPI-контракта
gen:
	go tool oapi-codegen -config oapi-codegen.yaml ../typespec/output/openapi.yaml

run:
	go run ./cmd/server

build:
	go build -o server ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy
