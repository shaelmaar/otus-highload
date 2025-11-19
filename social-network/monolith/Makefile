.PHONY: generate-server lint generate-sql generate-async-api
.PHONY: run-local stop-local stop-local-clear-volumes
.PHONY: run-local-infra stop-local-infra stop-local-infra-clear-volumes
.PHONY: run-simple-db-replicas stop-simple-db-replicas stop-simple-db-replicas-clear-volumes
.PHONY: run-patroni-db-replicas stop-patroni-db-replicas stop-patroni-db-replicas-clear-volumes
.PHONY: run-sharded-mongo stop-sharded-mongo stop-sharded-mongo-clear-volumes
.PHONY: run-monitoring stop-monitoring stop-monitoring-clear-volumes

build-docker-oapi:
	@docker build -f Dockerfile-oapi -t otus-highload-oapi-gen .

build-docker-gowrap:
	@docker build -f Dockerfile-gowrap -t otus-highload-gowrap .

build-docker-protoc-grpc:
	@docker build -f Dockerfile-protoc-grpc -t otus-highload-protoc-grpc .

generate-server: build-docker-oapi
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

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.6.1-alpine golangci-lint run --build-tags go_tarantool_ssl_disable --timeout 5m ./...

generate-sql: build-docker-gowrap
	@docker run --rm -v ${PWD}:${PWD} -w ${PWD} sqlc/sqlc:1.28.0 generate
	@docker run --rm \
    		-v $(PWD):/app \
    		-w /app \
    		otus-highload-gowrap gen -p ./internal/queries/pg -i Querier -o ./internal/queries/pg/querier_tx.go -t ./postgresql/dbtx.tmpl

generate-async-api:
	@docker run --rm -v $(PWD):/app -w /app lerenn/asyncapi-codegen asyncapi-codegen -i ./docs/asyncapi/asyncapi.json -g types -p async -o ./gen/async/types.gen.go

run-local:
	cd ./build/local && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-local:
	cd ./build/local && docker-compose --env-file .env.compose down --remove-orphans

stop-local-clear-volumes:
	cd ./build/local && docker-compose --env-file .env.compose down --remove-orphans -v

run-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-local-infra:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans

stop-local-infra-clear-volumes:
	cd ./build/local-infra && docker-compose --env-file .env.compose down --remove-orphans -v

run-simple-db-replicas:
	cd ./build/simple_db_replicas && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-simple-db-replicas:
	cd ./build/simple_db_replicas && docker-compose --env-file .env.compose down --remove-orphans

stop-simple-db-replicas-clear-volumes:
	cd ./build/simple_db_replicas && docker-compose --env-file .env.compose down --remove-orphans -v

run-patroni-db-replicas:
	cd ./build/patroni_db_replicas && docker-compose --env-file .env.compose up --build -d --force-recreate

stop-patroni-db-replicas:
	cd ./build/patroni_db_replicas && docker-compose --env-file .env.compose down --remove-orphans

stop-patroni-db-replicas-clear-volumes:
	cd ./build/patroni_db_replicas && docker-compose --env-file .env.compose down --remove-orphans -v

run-sharded-mongo:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose up  --build -d --force-recreate

rerun-sharded-mongo:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose up -d

stop-sharded-mongo:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose down --remove-orphans

stop-sharded-mongo-clear-volumes:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose down --remove-orphans -v

run-monitoring:
	cd ./build/monitoring && docker-compose up --build -d --force-recreate

stop-monitoring:
	cd ./build/monitoring && docker-compose down --remove-orphans

stop-monitoring-clear-volumes:
	cd ./build/monitoring && docker-compose down --remove-orphans -v

run-rabbitmq-cluster:
	cd ./build/rabbitmq-cluster && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-rabbitmq-cluster:
	cd ./build/rabbitmq-cluster && docker-compose --env-file .env.compose down --remove-orphans -v