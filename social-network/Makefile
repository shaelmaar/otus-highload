.PHONY: create-local-docker-net run-local stop-local stop-local-clear-volumes

create-local-docker-net:
	docker network create social-network-shared-net

run-local:
	cd ./build/local && docker-compose -f docker-compose-monolith.yaml -f docker-compose-dialogs.yaml --env-file .env.compose up  --build -d --force-recreate

stop-local:
	cd ./build/local && docker-compose -f docker-compose-monolith.yaml -f docker-compose-dialogs.yaml --env-file .env.compose down --remove-orphans

stop-local-clear-volumes:
	cd ./build/local && docker-compose -f docker-compose-monolith.yaml -f docker-compose-dialogs.yaml --env-file .env.compose down --remove-orphans -v
