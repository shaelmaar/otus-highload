## social-network
Заготовка соц сети
- monolith - монолит, ядро соц сети 
- dialogs - сервис диалогов

## Запуск сервисов


1. Скопируйте файл окружения:
    ```bash
   cp ./build/local/.env.monolith.example ./build/local/.env.monolith
   cp ./build/local/.env.dialogs.example ./build/local/.env.dialogs
   cp ./build/local/.env.compose.example ./build/local/.env.compose

2. При необходимости поменять значения переменных в .env и .env.compose
3. Создать общую сеть для докер сервисов:
    ```bash
   make create-local-docker-net

4. Запуск:
   ```bash
   make run-local

5. Остановка:
   ```bash
   make stop-local
   ```
   или с удалением всех хранилищ:
   ```bash
   make stop-local-clear-volumes
 