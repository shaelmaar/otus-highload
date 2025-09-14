## Запуск проекта

### Базовый запуск проекта:

1. Скопируйте файл окружения:
   ```bash
   cp .env.example .env.local
   
2. При необходимости поменять значения переменных в .env.local.

3. Запустите проект:
   ```bash
   make run
   
4. Чтобы остановить проект
   ```bash
   make stop
   
### Запуск с репликами бд (1 мастер, 2 слейва)

1. При необходимости поменять значения переменных в .env.local, .env.compose

2. Для запуска: 
   ```bash
   make run-simple-db-replicas
3. Для остановки:
   ```bash
   make stop-simple-db-replicas

Необходимые для билда файлы в ./build/simple_db_replicas


### Запуск с репликами бд с patroni (1 мастер, 2 слейва)

1. При необходимости поменять значения переменных в .env.local, .env.compose

2. Для запуска:
   ```bash
   make run-patroni-db-replicas
3. Для остановки:
   ```bash
   make stop-patroni-db-replicas

Необходимые для билда файлы в ./build/patroni_db_replicas