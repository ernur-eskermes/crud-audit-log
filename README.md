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

SERVER_PORT=9000
```

Use `make run` to build&run project, `make lint` to check code with linter, `make migrate` to apply the migration scheme.