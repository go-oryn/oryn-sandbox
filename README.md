# oryn-sandbox

Oryn Go framework experiementation sandbox.

## Start

Start the application with:
```shell
make fresh && make logs
```

This will expose:

- [http://localhost:8888](http://localhost:8888): API endpoints
- [http://localhost:8889/health?verbose=true](http://localhost:8889/health?verbose=true): Health check endpoints
- [http://localhost:3000](http://localhost:3000): Grafana LGTM stack (log logs, traces, metrics)
- [localhost:3306](localhost:3306): MySQL database

## Usage

This repository provides a [Makefile](Makefile):

```shell
make up      # start the docker compose stack
make down    # stop the docker compose stack
make logs    # stream the docker compose stack logs
make fresh   # refresh the docker compose stack (with db reset and seeding)
make migrate # run db migrations
make seed    # run db seeds
make test    # run tests
make lint    # run linter
```