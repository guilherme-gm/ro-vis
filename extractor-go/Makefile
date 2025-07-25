# Development/utility scripts

# Install unlisted dependencies for development
install-dev:
	go install github.com/githubnemo/CompileDaemon@latest

# Run extractor in development mode
extractor-dev:
	CGO_ENABLED=1 GOARCH=386 go run -tags luaa cmd/dev/dev.go

# Run miner tool
miner:
	CGO_ENABLED=1 GOARCH=386 go run -tags luaa cmd/miner/miner.go

test-extractor:
	CGO_ENABLED=1 GOARCH=386 go test -tags luaa -v ./...

# Applies database migrations
migrate:
	go run ./cmd/migrate/migrate.go

# Run the API in development mode, with auto restart.
# Requires CompileDaemon -- https://github.com/githubnemo/CompileDaemon
# 	go install github.com/githubnemo/CompileDaemon@latest
#
api-dev:
	CompileDaemon -directory="." -build="go build ./cmd/api/api.go" -command="./api"
