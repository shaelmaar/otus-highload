GEN_DIR := ./gen

.PHONY: generate-server lint generate-sql

build-docker-oapi:
	docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

build-docker-gowrap:
	docker build -f Dockerfile-gowrap -t otus-highload-gowrap .

generate-server: build-docker-oapi
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-oapi-gen -config oapi-codegen.yaml ./docs/openapi/swagger.yaml > ./gen/serverhttp/server.go

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.2.1-alpine golangci-lint run --timeout 5m ./...

generate-sql: build-docker-gowrap
	docker run --rm -v ${PWD}:${PWD} -w ${PWD} sqlc/sqlc:1.28.0 generate
	docker run --rm \
    		-v $(PWD):/app \
    		-w /app \
    		otus-highload-gowrap gen -p ./internal/queries/pg -i Querier -o ./internal/queries/pg/querier_tx.go -t ./postgresql/dbtx.tmpl