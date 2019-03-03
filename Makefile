# Include the .env file and export all environment variables. Swith to
# .env.development for development, or .env.production for production etc. The
# "-" symbol means that the .env file is optional - it will not throw error if
# the file is not found.
-include .env
export

# The git commit version allows us to track the version of the commit the
# application is using, so that we know that we are always deploying the latest
# version. The docker build will also be tagged with this version.
TAG := $(shell git rev-parse --short HEAD)

# The build date allows us to know when the service was last deployed, and how
# long have the service been running in production (uptime).
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

install: # Install required modules.
	@go get ./...
	GO111MODULE=on go get ./...

mod: # Initialize go modules and update dependencies.
	GO111MODULE=on go mod init
	GO111MODULE=on go mod tidy 
	GO111MODULE=on go get

start: # Start the main application with the exported environment variables.
	@go run main.go
