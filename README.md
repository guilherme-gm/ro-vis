# Installation

```
sudo apt-get install gcc-multilib
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/vektra/mockery/v3@v3.5.5
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Regenerating mocks
```
mockery
```

## Regenerating SQL
```
sqlc generate
```

## Testing

(It must run in x86 32bits)

```
make test-extractor
```
