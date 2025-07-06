GEN_DIR := ./gen

.PHONY: generate-server lint

build-docker-oapi:
	docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

generate-server: build-docker-oapi
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-oapi-gen -config oapi-codegen.yaml ./docs/openapi/swagger.yaml > ./gen/serverhttp/server.go

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.2.1-alpine golangci-lint run --timeout 5m ./...