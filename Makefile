# Include the .env file and export all environment variables. Swith to
# .env.development for development, or .env.production for production etc. The
# "-" symbol means that the .env file is optional - it will not throw error if
# the file is not found.
-include .env
ifeq ($(MAKE_ENV),)
	MAKE_ENV := development
endif
# This will throw an error if the file is not found.
include .env.$(MAKE_ENV)
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
	@go get -u github.com/pressly/goose/cmd/goose
	GO111MODULE=on go get ./...
	GO111MODULE=on go mod tidy 

mod: # Initialize go modules and update dependencies.
	GO111MODULE=on go mod init
	GO111MODULE=on go get ./...
	GO111MODULE=on go mod tidy 

start: # Start the main application with the exported environment variables.
	@go run main.go


DESCRIPTION := "go-microservice sample app"
NAME := $(shell git config --get user.name)/$(shell basename `git remote get-url origin` .git)
URL := "the url of the service"
CMD := "docker run -d ${NAME}"
VCS_REF := $(shell git rev-parse HEAD) 
VCS_URL := $(shell git remote get-url origin)
VENDOR := $(shell git config --get user.name)
VERSION := $(shell git rev-parse --short HEAD)

docker: 
	@docker-compose build app
	@docker tag ${NAME}:latest ${NAME}:${TAG}

inspect: 
	# Brew install jq to pretty print the json.
	@docker inspect --format='{{json .Config.Labels}}' ${NAME} | jq

up:
	@docker-compose up -d db

down:
	@docker-compose down


## DB
MIGRATION_FOLDER := "./migrations"
create-migration-%:
	@mkdir -p ${MIGRATION_FOLDER}
	@goose -dir ${MIGRATION_FOLDER} create $* sql 

migrate:
	@goose -dir ${MIGRATION_FOLDER} mysql "${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}?parseTime=true" up

rollback:
	@goose -dir ${MIGRATION_FOLDER} mysql "${DB_USER}:${DB_PASS}@${DB_HOST}/${DB_NAME}?parseTime=true" down

clean:
	@rm -rf tmp

mysql:
	@mysql -h ${DB_HOST} -u ${DB_USER} -p ${DB_NAME}
