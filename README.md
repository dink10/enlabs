# Payment service #

#### Service for processing incoming requests with payment

## Dependencies

- golang
- docker and docker-compose
- swagger (https://github.com/swaggo/swag)

golangci-lint (https://github.com/golangci/golangci-lint) if you want to run linter.

## How to run

- make or `make up` builds and runs application in docker containers. It uses `deployment/docker-compose.yml`
to start application, database and migrations.
- `make build` builds application in docker containers.
- `make start` starts application in docker containers.
- `make stop` starts application in docker containers.

`make lint` to check if your code passes linter.

## How to use

1. Make sure ports :8085/:5432 are free
2. Run `make up`
3. Run tests `make test`
4. See swagger documentation below

### Swagger documentation

After running make up/make start open http://localhost:8085/swagger/index.html in your browser

### Environment variables

HTTP server configuration:
```
SERVER_HOST: 0.0.0.0
SERVER_PORT: 3000
SERVER_LOG_REQUESTS: 1
```

Database configuration:
```
DB_HOST: postgres
DB_PORT: 5432
DB_NAME: postgres
DB_USER: admin
DB_PASSWORD: admin
DB_MAX_CONN: 10
DB_ENABLE_LOG: 1
```

Logger configuration:
```
LOG_LEVEL: debug
```