## Запуск проекта

### Запуск со всеми зависимостями:

1. Скопируйте файл окружения:
    ```bash
   cp ./build/local/.env.example ./build/local/.env
   cp ./build/local/.env.compose.example ./build/local/.env.compose

2. При необходимости поменять значения переменных в .env и .env.compose

3. Запуск:
   ```bash
   make run-local
   
4. Остановка:
   ```bash
   make stop-local
   ```
   или с удалением всех хранилищ:
   ```bash
   make stop-local-clear-volumes

### Запуск инфраструктуры
#### Для дебага, запуска из IDE итп

1. Скопируйте файлы окружения:
   ```bash
   cp ./build/local-infra/.env.example ./build/local-infra/.env
   cp ./build/local-infra/.env.compose.example ./build/local-infra/.env.compose

2. При необходимости поменять значения переменных в .env и .env.compose
3. Запуск инфраструктуры:
   ```bash
   make run-local-infra
4. Прокинуть в приложение .env.
5. Остановка:
   ```bash
   make stop-local-infra
   ```
   или с удалением всех хранилищ:
   ```bash
   make stop-local-infra-clear-volumes
   ```
   
### Запуск с репликами бд (1 мастер, 2 слейва)

1. Скопируйте файлы окружения:
   ```bash
   cp ./build/simple_db_replicas/.env.example ./build/simple_db_replicas/.env
   cp ./build/simple_db_replicas/.env.compose.example ./build/simple_db_replicas/.env.compose
   
2. Поменять значения если нужно

3. Для запуска: 
   ```bash
   make run-simple-db-replicas
4. Для остановки:
   ```bash
   make stop-simple-db-replicas
   ```
   или с удалением всех хранилищ:
   ```bash
   make stop-simple-db-replicas-clear-volumes
   ```

Необходимые для билда файлы в ./build/simple_db_replicas


### Запуск с репликами бд с patroni (1 мастер, 2 слейва)

1. Скопируйте файлы окружения:
   ```bash
   cp ./build/patroni_db_replicas/.env.example ./build/patroni_db_replicas/.env
   cp ./build/patroni_db_replicas/.env.compose.example ./build/patroni_db_replicas/.env.compose
   
2. Поменять значения если нужно

3. Для запуска:
   ```bash
   make run-patroni-db-replicas
4. Для остановки:
   ```bash
   make stop-patroni-db-replicas
   ```
   или с удалением всех хранилищ:
   ```bash
   make stop-patroni-db-replicas-clear-volumes
   ```
   
Необходимые для билда файлы в ./build/patroni_db_replicas

### Запуск мониторинга контейнеров

1. Поменять конфигурацию в ./build/monitoring/prometheus.yaml при необходимости
2. Для запуска:
   ```bash
   make run-monitoring
3. После успешного запуска импортировать dashboard
4. Для остановки:
    ```bash
   make stop-monitoring
   ```
   или с удалением всех хранилищ:
   ```bash
   make stop-monitoring-clear-volumes
   ```
   

## Для разработки

1. Генерация http-сервера на основе спецификации OpenAPI 3, лежащей в docs/openapi/swagger.yaml:
   ```bash
   make generate-http-server
   ```
2. Генерация grpc-сервера на основе protobuf контрактов, лежащих в docs/protobuf/grpc/server:
   ```bash
   make generate-grpc-server
   ```
3. Генерация grpc-клиента dialogs сервиса на основе protobuf контрактов, лежащих в docs/protobuf/grpc/dialogs:
   ```bash
   make generate-grpc-dialogs-client
   ```
4. Генерация sql из схемы бд (postgresql/migrations) и запросов (postgresql/queries):
   ```bash
   make generate-sql
   ```
5. Линтер golang-кода:
   ```bash
   make lint
   ```
