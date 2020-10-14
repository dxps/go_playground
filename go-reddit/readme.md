## GoReddit

Supposely a Reddit-like app project.

<br/>

### Packages

The project is based on a couple of commonly used packages such as [`sqlx`](github.com/jmoiron/sqlx) and [`pq`](github.com/lib/pq). Of course, for all the included dependencies, check out `go.mod` file.

<br/>

### Database Migrations

The CLI (binary) release of [golang-migrate/migrate](github.com/golang-migrate/migrate) is used: downloaded and installed into `${GOPATH}/bin` as `go_migrate`.

Initial (and any later introduced elements part of the ) data model has been provisioned using `make db-migrate`.
