.PHONY: run stop run-simple-db-replicas stop-simple-db-replicas generate-server lint generate-sql

run:
	docker-compose --env-file .env.compose up  --build -d --force-recreate

stop:
	docker-compose down --remove-orphans

run-simple-db-replicas:
	docker-compose -f docker-compose_simple-db-replicas.yaml --env-file .env.compose up  --build -d --force-recreate

stop-simple-db-replicas:
	docker-compose -f docker-compose_simple-db-replicas.yaml down --remove-orphans

build-docker-oapi:
	@docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

build-docker-gowrap:
	@docker build -f Dockerfile-gowrap -t otus-highload-gowrap .

generate-server: build-docker-oapi
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-oapi-gen -config oapi-codegen.yaml ./docs/openapi/swagger.yaml > ./gen/serverhttp/server.go

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.2.1-alpine golangci-lint run --timeout 5m ./...

generate-sql: build-docker-gowrap
	@docker run --rm -v ${PWD}:${PWD} -w ${PWD} sqlc/sqlc:1.28.0 generate
	@docker run --rm \
    		-v $(PWD):/app \
    		-w /app \
    		otus-highload-gowrap gen -p ./internal/queries/pg -i Querier -o ./internal/queries/pg/querier_tx.go -t ./postgresql/dbtx.tmpl