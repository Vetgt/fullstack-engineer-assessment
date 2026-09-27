# Fullstack Engineer Assessment - Backend

Backend implementation for the task-management assessment.

## Stack

- Go 1.27.x
- Gin 1.12
- MySQL 8.4
- Redis 8
- `database/sql` + MySQL driver
- go-redis v9

## What is implemented

- `GET /api/tasks` filtering by `status`, `keyword`, `assignee`
- pagination with `page` and `limit`
- sorting with a whitelist through `sort`
- `GET /api/tasks/:id`
- `POST /api/tasks`
- `PUT /api/tasks/:id`
- `DELETE /api/tasks/:id` as soft delete
- duplicate title returns HTTP 409
- consistent error body
- Redis list cache with 60-second TTL
- cache key includes query parameters
- cache invalidation after create/update/delete
- soft-deleted tasks are excluded from list/get/search
- SQL migration
- seed command
- unit tests for update, search, and cache invalidation

## First run

1. Copy `.env.example` to `.env`.
2. Start dependencies:

```bash
docker compose up -d
```

3. Download Go modules:

```bash
go mod tidy
```

4. Run migration:

```bash
go run ./cmd/migrate
```

5. Insert sample data:

```bash
go run ./cmd/seed
```

6. Start the API:

```bash
go run ./cmd/server
```

7. Check health:

```bash
curl http://localhost:8080/health
```

## API examples

List:

```bash
curl "http://localhost:8080/api/tasks?page=1&limit=10&keyword=login&status=todo&assignee=Betran&sort=created_at_desc"
```

Create:

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Prepare assessment","description":"Finish fullstack assessment","status":"todo","assignee":"Betran"}'
```

Update:

```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Fix login validation","description":"Updated description","status":"done","assignee":"Betran"}'
```

Soft delete:

```bash
curl -i -X DELETE http://localhost:8080/api/tasks/1
```

## Test

```bash
go test ./...
```

## Cache behavior

Every list request uses a Redis key that contains its query string. Example:

```text
tasks:list:assignee=Betran&keyword=login&limit=10&page=1&sort=created_at_desc&status=todo
```

Successful create/update/delete calls invalidate all `tasks:list:*` keys. A successful list response also exposes `X-Cache: HIT` or `X-Cache: MISS`.
