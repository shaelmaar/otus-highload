.PHONY: generate-http-server generate-grpc-server generate-grpc-dialogs-client lint lint-clean-cache generate-sql generate-async-api
.PHONY: run-local-infra stop-local-infra stop-local-infra-clear-volumes

build-docker-oapi:
	@docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

build-docker-gowrap:
	@docker build -f Dockerfile-gowrap -t otus-highload-gowrap .

generate-http-server: build-docker-oapi
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-oapi-gen -config oapi-codegen.yaml ./docs/openapi/swagger.yaml > ./gen/serverhttp/server.go

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.7.1-alpine golangci-lint run --build-tags go_tarantool_ssl_disable --timeout 5m ./...

lint-clean-cache:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.6.1-alpine golangci-lint cache clean

generate-sql: build-docker-gowrap
	@docker run --rm -v ${PWD}:${PWD} -w ${PWD} sqlc/sqlc:1.28.0 generate
	@docker run --rm \
    		-v $(PWD):/app \
    		-w /app \
    		otus-highload-gowrap gen -p ./internal/queries/pg -i Querier -o ./internal/queries/pg/querier_tx.go -t ./postgresql/dbtx.tmpl

generate-async-api:
	@docker run --rm -v $(PWD):/app -w /app lerenn/asyncapi-codegen asyncapi-codegen -i ./docs/asyncapi/asyncapi.json -g types -p async -o ./gen/async/types.gen.go

run-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans

stop-local-infra-clear-volumes:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans -v
