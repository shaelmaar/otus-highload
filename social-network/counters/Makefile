.PHONY: generate-http-server lint
.PHONY: run-local stop-local stop-local-clear-volumes

build-docker-oapi:
	@docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

build-docker-protoc-grpc:
	@docker build -f Dockerfile-protoc -t otus-highload-protoc .

generate-http-server: build-docker-oapi
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-oapi-gen -config oapi-codegen.yaml ./docs/openapi/swagger.yaml > ./gen/serverhttp/server.go

generate-kafka-producer: build-docker-protoc-grpc
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-protoc --go_out=gen/kafka/producer/proto \
                             		--go_opt=paths=source_relative \
                             		--proto_path=./docs/protobuf/kafka/producer \
                             		./docs/protobuf/kafka/producer/*.proto

generate-kafka-consumer: build-docker-protoc-grpc
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		otus-highload-protoc --go_out=gen/kafka/consumer/proto \
                             		--go_opt=paths=source_relative \
                             		--proto_path=./docs/protobuf/kafka/consumer \
                             		./docs/protobuf/kafka/consumer/*.proto

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.7.1-alpine golangci-lint run --timeout 5m ./...

run-local:
	cd ./build/local && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-local:
	cd ./build/local && docker-compose --env-file .env.compose down --remove-orphans

stop-local-clear-volumes:
	cd ./build/local && docker-compose --env-file .env.compose down --remove-orphans -v


