# Search Filters API

Минимальный сервис для хранения и применения сохранённых поисковых фильтров.

## Возможности
- Создание фильтров (`POST /filters`)
- Получение списка фильтров (`GET /filters`)
- Получение фильтра по ID (`GET /filters/{id}`)
- Обновление фильтра (`PUT /filters/{id}`)
- Удаление фильтра (`DELETE /filters/{id}`)
- Применение фильтра с подстановкой плейсхолдеров (`GET /filters/{id}/apply`)

## Динамические плейсхолдеры
- `{{today}}` → текущая дата (UTC).
- `{{today-7d}}` → дата 7 дней назад.
- `{{current_user}}` → ID пользователя.

## Установка и запуск

### Требования
- Go 1.22+
- PostgreSQL

### Настройки
Конфигурация задаётся через YAML-файл (`CONFIG_FILE`) и переменные окружения:

```yaml
timezone: "Europe/Moscow"
postgres_host: "localhost"
postgres_port: "5432"
postgres_db: "searchfilt"
```

А также переменные окружения:

```bash
export POSTGRES_USER=<username>
export POSTGRES_PASSWORD=<password>
```

### Makefile команды
В проекте есть удобный `Makefile`:

```makefile
.PHONY: run migrate-up migrate-down tidy

run:
	go run ./cmd/app serve

migrate-up:
	go run ./cmd/app migrate up

migrate-down:
	go run ./cmd/app migrate down

tidy:
	go mod tidy
```

### Примеры запросов

Создать фильтр:
```bash
curl -s -X POST http://localhost:8080/filters \
  -H "Content-Type: application/json" \
  -d '{"name":"Go articles","query":{"tags":["golang"],"date_from":"{{today-7d}}"}}'
```

Получить список:
```bash
curl -s http://localhost:8080/filters | jq
```

Получить фильтр:
```bash
curl -s http://localhost:8080/filters/1 | jq
```

Обновить фильтр:
```bash
curl -s -X PUT http://localhost:8080/filters/1 \
  -H "Content-Type: application/json" \
  -d '{"query":{"tags":["golang","db"],"date_from":"2025-07-27"}}' | jq
```

Применить фильтр:
```bash
curl -s http://localhost:8080/filters/1/apply | jq
```

Удалить фильтр:
```bash
curl -i -X DELETE http://localhost:8080/filters/1
```

---
