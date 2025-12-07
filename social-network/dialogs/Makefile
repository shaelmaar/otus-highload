.PHONY: generate-http-server generate-grpc-server generate-grpc-monolith-client lint
.PHONY: run-local-infra stop-local-infra stop-local-infra-clear-volumes

build-docker-oapi:
	@docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

build-docker-protoc-grpc:
	@docker build -f Dockerfile-protoc-grpc -t otus-highload-protoc-grpc .

generate-http-server: build-docker-oapi
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-oapi-gen -config oapi-codegen.yaml ./docs/openapi/swagger.yaml > ./gen/serverhttp/server.go

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.7.1-alpine golangci-lint run --build-tags go_tarantool_ssl_disable --timeout 5m ./...

run-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans

stop-local-infra-clear-volumes:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans -v
