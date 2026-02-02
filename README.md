# HR Engine
=========

HR Engine is a Go library providing helper utilities and a lightweight HTTP server
structure for creating api with ease. The project includes common
helpers for configuration, database adapters (Postgres/SQLite), pagination helpers, and
web request/response scaffolding.

## Features
--------
- Configuration loading and environment helpers
- Database abstraction layer with Postgres and SQLite implementations
- HTTP server scaffolding and request/response helpers
- Common utilities and error handling helpers

## Installing
-----------
```
go get github.com/rafirandoni/hrengine
```

## Project Layout
--------------
- `hrengine.go` — package entry / top-level wiring
- `common/` — shared helpers (pagination, etc.)
- `constant/` — project constants
- `db/` — database adapters (`db_pgsql.go`, `db_sqlite.go`, `db.go`)
- `repository/` — repository interfaces and implementations
- `server/` — HTTP server and routing (`server.go`, `http.go`)
- `utility/` — environment, error helpers, utilities
- `web/` — request/response helpers and web glue

## Contributing
------------
Contributions are welcome. Please open issues for bugs or feature requests and submit
pull requests for proposed changes. Keep changes focused and provide tests where
appropriate.

## License
-------
Check the repository for a license file or add one as appropriate for your project.

## Contact
-------
For questions about this library, open an issue in the repository.
