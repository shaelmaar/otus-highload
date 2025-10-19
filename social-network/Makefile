.PHONY: run stop run-local-infra stop-local-infra generate-server lint generate-sql
.PHONY: run-simple-db-replicas stop-simple-db-replicas stop-simple-db-replicas-clear-volumes
.PHONY: run-patroni-db-replicas stop-patroni-db-replicas run-monitoring stop-monitoring
.PHONY: run-sharded-mongo stop-sharded-mongo stop-sharded-mongo-clear-volumes
.PHONY: generate-server lint generate-sql

run:
	docker-compose --env-file .env.compose up  --build -d --force-recreate

stop:
	docker-compose --env-file .env.compose down --remove-orphans

run-local-infra:
	docker-compose -f docker-compose-local.yaml --env-file .env.compose up  --build -d --force-recreate

stop-local-infra:
	docker-compose -f docker-compose-local.yaml --env-file .env.compose down --remove-orphans

run-simple-db-replicas:
	cd ./build/simple_db_replicas && docker-compose --env-file .env.compose up  --build -d --force-recreate

stop-simple-db-replicas:
	cd ./build/simple_db_replicas && docker-compose --env-file .env.compose down --remove-orphans

stop-simple-db-replicas-clear-volumes:
	cd ./build/simple_db_replicas && docker-compose --env-file .env.compose down --remove-orphans -v

run-sharded-mongo:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose up  --build -d --force-recreate

rerun-sharded-mongo:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose up -d

stop-sharded-mongo:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose down --remove-orphans

stop-sharded-mongo-clear-volumes:
	cd ./build/sharded-mongo && docker-compose --env-file .env.compose down --remove-orphans -v

run-patroni-db-replicas:
	docker-compose -f docker-compose_patroni-db-replicas.yaml --env-file .env.compose up  --build -d --force-recreate

stop-patroni-db-replicas:
	docker-compose -f docker-compose_patroni-db-replicas.yaml --env-file .env.compose down --remove-orphans -v

run-monitoring:
	docker-compose -f docker-compose_mon.yaml  up  --build -d --force-recreate

stop-monitoring:
	docker-compose -f docker-compose_mon.yaml  down --remove-orphans

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