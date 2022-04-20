# Сервис логирования событий для CRUD APP
## Пример приложения из недели #7

### Стэк
- go 1.17
- postgres
- [golangci-lint](https://github.com/golangci/golangci-lint) (<i>optional</i>, used to run code checks)

Create .env file in root directory and add following values:
```dotenv
MONGO_URI=mongodb://audit-db
MONGO_DATABASE=audit
MONGO_USER=admin
MONGO_PASSWORD=g0langn1nja

GRPC_PORT=9000

HTTP_HOST=localhost
HTTP_PORT=9001
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=10s
HTTP_MAX_HEADER_MEGABYTES=1
```

Use `make run` to build&run project, `make lint` to check code with linter, `make migrate` to apply the migration scheme.
### Generating .pb.go files
```protoc --go_out=proto --go-grpc_out=proto proto/audit.proto```
