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

PROJECTNAME := $(shell basename $(PWD))

.DEFAULT_GOAL=help
.PHONY: help
help:
	@echo 
	@echo "Run a command for ${PROJECTNAME}:"
	@echo 
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[92m%s:\033[0m %s\n", $$1, $$NF}' $(MAKEFILE_LIST) | column -t -s ":" | sort
	@echo 

# The git commit version allows us to track the version of the commit the
# application is using, so that we know that we are always deploying the latest
# version. The docker build will also be tagged with this version.
TAG := $(shell git rev-parse --short HEAD)

# The build date allows us to know when the service was last deployed, and how
# long have the service been running in production (uptime).
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

MODULES := go.mod go.sum

install: $(MODULES) ## install required modules
	@go get ./...
	@go get -u github.com/pressly/goose/cmd/goose
	@#The release tag is not updated, need to point to master to ensure it is pulling the latest version.
	@GO111MODULE=on go get github.com/satori/go.uuid@master
	@GO111MODULE=on go get
	@GO111MODULE=on go mod tidy 

$(MODULES): 
	# GO111MODULE=on go mod init
	GO111MODULE=on go get ./...
	GO111MODULE=on go mod tidy 

.PHONY: start test
start: up ## runs the main application 
	@# @gin -p ${PORT} run main.go
	@go run main.go

test: ## test the application
	go test -v ./...

DESCRIPTION := "go-microservice sample app"
NAME := $(shell git config --get user.name)/$(shell basename `git remote get-url origin` .git)
URL := "the url of the service"
CMD := "docker run -d ${NAME}"
VCS_REF := $(shell git rev-parse HEAD) 
VCS_URL := $(shell git remote get-url origin)
VENDOR := $(shell git config --get user.name)
VERSION := $(shell git rev-parse --short HEAD)

.PHONY: docker inspect up down

docker: ## build the docker image
	@docker-compose build app
	@docker tag ${NAME}:latest ${NAME}:${TAG}

inspect: ## inspect the docker labels
	# Brew install jq to pretty print the json.
	@docker inspect --format='{{json .Config.Labels}}' ${NAME} | jq

up: ## bring docker-compose up
	@docker-compose up -d db

down: ## bring docker-compose down
	@docker-compose down

MIGRATION_FOLDER := "./migrations"
.PHONY: migrate rollback clean mysql

create-migration-%: ## creates a new migration file
	@mkdir -p ${MIGRATION_FOLDER}
	@goose -dir ${MIGRATION_FOLDER} create $* sql

migrate: ## run the migrations to the latest
	@goose -dir ${MIGRATION_FOLDER} mysql "${DB_USER}:${DB_PASS}@tcp(${DB_HOST})/${DB_NAME}?parseTime=true" up

rollback: ## rollback a migration version
	@goose -dir ${MIGRATION_FOLDER} mysql "${DB_USER}:${DB_PASS}@tcp(${DB_HOST})/${DB_NAME}?parseTime=true" down

clean: ## clean the temporary database directory
	@rm -rf tmp

mysql: ## access the mysql cli
	@mysql -h ${DB_HOST} -u ${DB_USER} -p ${DB_NAME}
