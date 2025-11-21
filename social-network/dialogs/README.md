## Запуск проекта

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


## Для разработки

1. Генерация http-сервера на основе спецификации OpenAPI 3, лежащей в docs/openapi/swagger.yaml:
   ```bash
   make generate-http-server
   ```
2. Генерация grpc-сервера на основе protobuf контрактов, лежащих в docs/protobuf/grpc/server:
   ```bash
   make generate-grpc-server
   ```
3. Генерация grpc-клиента monolith сервиса на основе protobuf контрактов, лежащих в docs/protobuf/grpc/monolith:
   ```bash
   make generate-grpc-monolith-client
   ```
4. Линтер golang-кода:
   ```bash
   make lint
   ```
