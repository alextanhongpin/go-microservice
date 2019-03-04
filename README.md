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


## Config

TL;DR;

- there are two types config - `global` and `package`. 
- `global` config includes app specific configuration, e.g. `APP_PORT`, `APP_HOST`, `APP_VERSION`, `APP_BUILD_AT`
- `package` config are configuration for vendor packages, such as database, logger etc. `DB_NAME`, `DB_HOST`, `DB_PASS`, `DB_USER`
- configs can have sane defaults for the `development`, `production`, or `nop` (null object pattern) 
- configs could be passed through golang `flag` or `envvar` (environment variables), pick one and standardize it
- include the `.env` in the `.gitignore`, we do not want to commit sensitive info to git repository
