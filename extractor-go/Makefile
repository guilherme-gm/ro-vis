# Development/utility scripts

# Install unlisted dependencies for development
install-dev:
	go install github.com/githubnemo/CompileDaemon@latest

# Run the API in development mode, with auto restart.
# Requires CompileDaemon -- https://github.com/githubnemo/CompileDaemon
# 	go install github.com/githubnemo/CompileDaemon@latest
#
api-dev:
	CompileDaemon -directory="." -build="go build ./cmd/api/api.go" -command="./api"
