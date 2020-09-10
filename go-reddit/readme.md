## GoReddit

Supposely a Reddit-like app project.

### Packages

The following packages are used by the implementation:

- `sqlx`
- _to-be-cont'd_

### Database Migrations

The CLI (binary) release of [golang-migrate/migrate](github.com/golang-migrate/migrate) is used: downloaded and installed into \$GOPATH/bin as `go_migrate`.

Initial data model has been provisioned using `make db-migrate`.
