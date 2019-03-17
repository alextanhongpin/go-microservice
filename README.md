# go-microservice

Some thoughts on designing maintainable microservice with golang.

- background worker
- [config](#config)
- controller, model, service and repository
- database (migration etc)
- dependency management
- documentation
- error handling
- goroutines termination
- [graceful shutdown](#graceful-shutdown)
- [health](#health)
- internal service call
- [logging](#logging)
- middlewares (authentication, roles and scopes, invalidating tokens)
- minimal docker build
- naming convention
- pkg vs model
- [request id](#request-id)
- testing
- validation

## Commit

The commit messages are based on [Semantic Commit Messages](https://seesparkbox.com/foundry/semantic_commit_messages).

## Setup

Learn to use a `Makefile`, it simplifies command and standardize your development workflow. Here are some commonly used command:

```bash
# Initialize the project if it is not yet initialized.
$ make init

# Install all the dependencies required for this project.
$ make install

# Start a local development server with the environment variables set.
# It doesn't matter if you are starting a web project or backend server, you can standardize the command to start your app.
$ make start

# Build a binary, or compile a web project.
$ make build

# Build a docker image.
$ make docker

# Run test.
$ make test

# Stop the server.
$ make stop

# Start docker-compose locally (normally for running a development database).
$ make up

# Stop docker-compose locally.
$ make down

# Clean up temporary directory/resources that are used locally.
$ make clean
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
- pass the config down through DI (dependency injection) or params, **DO NOT** call it straight from `os.Getenv`

## Working with different environment

We will have one base `.env.development` environent file that exports all the required environment variables for development environment. To override part of the environment, say for staging, just run:

```bash
MAKE_ENV=staging make your-command
```

## Database

TODO:

- Prefer uuid over auto-incrementing id
- Store uuid as `Binary(16)`
- MySQL 8.0 and above has support for `uuid_to_bin(uuid(), true)` and `bin_to_uuid(uuid(), true)` functions. The second arguments is set to true, which will rearrange the time component of the uuid to enhance indexing performance (by ordering it chronologically). This only works for uuid v1.
- MySQL uses uuid v1. If you are using a golang library to create the uuid externally, make sure the uuid used is the v1 version.
- paging with cursor pagination
- migrations files and execution

## Request ID

TODO:

- reasons to have
- usage example
- lifecycle - creation to received
- tracing with request id

References:

- https://stackoverflow.com/questions/25433258/what-is-the-x-request-id-http-header
- https://blog.heroku.com/http_request_id_s_improve_visibility_across_the_application_stack

## Graceful Shutdown

TODO

## Logging

- debugging in development
- format in production
- centralized logging
- removal of logs from certain endpoints, e.g. `/health`
- value to noise signal ratio - not all logs are good. know what to log
- log the request whenever there are errors - this allows us to trace which requests are causing the error. But remember not to log sensitive requests (passwords etc)
- wrap the errors and print out the stack trace whenever an error occurred

## Health Endpoint

Useful application metrics includes:

- the git commit version - allows us to know what is the latest version of the application deployed
- uptime - how long has the application been running before restarting
- deployed_at - when was the application deployed (or when the docker image is built)
