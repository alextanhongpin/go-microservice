# go-microservice

Some thoughts on designing maintainable microservice with golang.

- background worker
- config
- controller, model, service and repository
- database (migration etc)
- dependency management
- documentation
- error handling
- goroutines termination
- graceful shutdown
- health endpoint
- internal service call
- logging
- middlewares (authentication, roles and scopes, invalidating tokens)
- minimal docker build
- naming convention
- pkg vs model
- request id
- testing
- validation

## Setup

Learn to use a `Makefile`, it simplifies command and standardize your development workflow. Here are some commonly used command:

- `make init`: initialize the project if it is not yet initialized
- `make install`: install all the dependencies required for this project
- `make start`: start a local development server with the environment variables set. It doesn't matter if you are starting a web project or backend server, you can standardize the command to start your app.
- `make build`: build a binary, or compile a web project
- `make docker`: build a docker image
- `make test`: run test
- `make stop`: stop the server
- `make up`: start docker-compose locally (normally for running a development database)
- `make down`: stop docker-compose locally
- `make clean`: clean up temporary directory/resources that are used locally

For example, to start this project, you just need to clone the git repository, then run:

```bash
$ make install
$ make start
```

## Config

TL;DR;

- there are two types config - `global` and `package`. 
- `global` config includes app specific configuration, e.g. `APP_PORT`, `APP_HOST`, `APP_VERSION`, `APP_BUILD_AT`
- `package` config are configuration for vendor packages, such as database, logger etc. `DB_NAME`, `DB_HOST`, `DB_PASS`, `DB_USER`
- configs can have sane defaults for the `development`, `production`, or `nop` (null object pattern) 
- configs could be passed through golang `flag` or `envvar` (environment variables), pick one and standardize it
- include the `.env` in the `.gitignore`, we do not want to commit sensitive info to git repository
- there many libraries to parse and read environment config, use the one that is the most simple to use
- pass the config down through DI (dependency injection) or params, __DO NOT__ call it straight from `os.Getenv`
