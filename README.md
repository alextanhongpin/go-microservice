# go-microservice

The application, a part of the whole [architecture](https://github.com/alextanhongpin/go-microservice-architecture)

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
- [validation](#validation)

## TODO:
implement this
- https://docs.microsoft.com/en-us/azure/architecture/patterns/


Reliability: 
- circuit breaker (can be part of sidecar): prevents thundering herd when requests fails
- throttle (can be part of sidecar): prevents DDOS
- has graceful shutdown: server terminates all running processes safely

Observability:
- request id: has a unique id propagated towards each component
- tracing: paths taken are covered and tested
- logging: logs request/response

Security:
- no tokens/secrets in code
- no secure data logged to stdout

Docker:
- has scripts to build image
- has dockerignore to prevent large builds

Documentation
- has documentation on how to create the layers service/controller/repository/entity
- has scripts to start the program locally
- has dockerfiles
- has database (and all related stuff to database)
- has pre commit to lint, fix and run tests

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

https://github.com/alextanhongpin/evolutionary-architecture/blob/master/configuration.md

## Working with different environment

We will have one base `.env.development` environent file that exports all the required environment variables for development environment. To override part of the environment, say for staging, just run:

```bash
MAKE_ENV=staging make your-command
```

- limit access to production environment (e.g. mysql dump/backup)

## Database

- Prefer uuid over auto-incrementing id
- Store uuid as `Binary(16)`
- MySQL 8.0 and above has support for `uuid_to_bin(uuid(), true)` and `bin_to_uuid(uuid(), true)` functions. The second arguments is set to true, which will rearrange the time component of the uuid to enhance indexing performance (by ordering it chronologically). This only works for uuid v1.
- MySQL uses uuid v1. If you are using a golang library to create the uuid externally, make sure the uuid used is the v1 version.
- paging with cursor pagination
- migrations files and execution
- use prepared statement for golang to check errors in statement quick

## Request ID


- use a middleware to generate a unique request id for every request
- pass the request id down through context, and make every function accepts context as the first argument
- log the request id whenever there's an error, or when an operation succeeds to trace the steps
- log the error with the request id so that you can trace the error steps back from the log

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
- using global logger is okay, since logging happens in all layers of the system - passing them down to every layer is cumbersome

## Validation

- using global validator is okay, since validation is part of the business logic and is something you won't need to mock (?)
- ensure the requests are validated before calling the service
- validation should happen in the business logic layer, not controller. This makes testing easier, since it will reduce the negative conditions (service can only be executed if the validation pass), and we can test the service directly. Testing the controller is not an option (probably for integration testing, but for most cases, we want to test the business logic and skip the middleware setups for auth etc that is executed with the controller).
- passing down struct is more convenient that individual arguments since we can validate the struct as whole rather than individual params. you can still pass down individual params, but have a struct within the service with a validation tag, assign those params to the struct and validate the struct.
- ensure all conditions are met before calling the service - required fields should be clearly defined
- trim the strings before checking the length to ensure empty strings with spaces is not passed in
- for numbers, ensure it cannot be negative for pagination etc, always set a min max too for pagination to avoid abuse

## UseCases

- create a new usecase e.g. usecase.login.go
- define the request/response pair
- implement the function
- create a factory for the use case
- test the use case independently
- add scenarios on the go
- the service struct should hold all the use cases
- combine usecases in the service (usecases with include statement)
- better to mock a behavior than a dependency (e.g. rather than mocking a jwt provider/dependency, it is better to create a struct with the provider and a method that calls the provider, then create an interface on top of it)

## Health Endpoint

Useful application metrics includes:

- the git commit version - allows us to know what is the latest version of the application deployed
- uptime - how long has the application been running before restarting
- deployed_at - when was the application deployed (or when the docker image is built)
- do not put the database connection ping here (?)

## Testing

- to make mocking easier, pass in the struct that needs to be mocked

```go

package main

import (
	"fmt"
)

type mockSigner struct {
	key          string
	token        string
	err          error
	invoked      bool
	invokedCount int
}

func (m *mockSigner) Sign(key string) (string, error) {
	defer func() {
		m.invoked = true
		m.invokedCount++
	}()
	if m.key == key {
		return m.token, nil
	}
	return m.token, m.err
}

type Signer interface {
	Sign(key string) (string, error)
}

func main() {
	m := mockSigner{
		key:   "test",
		token: "xyz",
	}
	res, _ := m.Sign("test")
	fmt.Println(res)
}
```

## Roles and Scopes

Roles and scopes limits the API access to certain users, whether it is authenticated or not. Each API will have it's own scope (grouped by usecases), and only certain roles can access it. 


## Thoughts

I realised I've been mixing the interface layers and their implementation in the `domain` folders - they should ideally be separated. The reason being that the `domain` layer is one that can be reused even when the implementations has been switched. Thus it should purely be defining the business rules, the `why`, not `how`.

Also, from Rakyll's [blog](https://rakyll.org/).
```
Naming patterns based on other languages’ dependency inversion conventions are anti-patterns in Go. Naming styles such the following don’t fit into the Go ecosystem.
```
```go
type Banana interface {
    //...
}
type BananaImpl struct {}
```

Makes we think twice about how to name interfaces. Will probably refactor the implementation again.
