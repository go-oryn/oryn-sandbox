.PHONY: up down fresh logs test lint

up:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose up -d

down:
	docker compose down

migrate:
	docker compose exec oryn-app go run . migrate up

seed:
	docker compose exec oryn-app go run . seed

fresh:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V
	docker compose exec oryn-app go run . migrate reset
	docker compose exec oryn-app go run . migrate up
	docker compose exec oryn-app go run . seed

logs:
	docker compose logs -f

test:
	go test -v -race -cover -count=1 -failfast ./...

lint:
	golangci-lint run -v