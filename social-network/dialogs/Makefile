.PHONY: generate-http-server generate-grpc-server generate-grpc-monolith-client lint
.PHONY: run-local stop-local stop-local-clear-volumes
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

generate-grpc-server: build-docker-protoc-grpc
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-protoc-grpc -I docs/protobuf/grpc/server \
		  --go_out=gen/servergrpc --go_opt=paths=source_relative \
		  --go-grpc_out=gen/servergrpc --go-grpc_opt=paths=source_relative \
		  docs/protobuf/grpc/server/*.proto

generate-grpc-monolith-client: build-docker-protoc-grpc
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-protoc-grpc -I docs/protobuf/grpc/monolith \
		  --go_out=gen/clientgrpc/monolith --go_opt=paths=source_relative \
		  --go-grpc_out=gen/clientgrpc/monolith --go-grpc_opt=paths=source_relative \
		  docs/protobuf/grpc/monolith/*.proto

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.6.1-alpine golangci-lint run --build-tags go_tarantool_ssl_disable --timeout 5m ./...


run-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans

stop-local-infra-clear-volumes:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans -v
