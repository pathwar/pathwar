COMPOSE = docker-compose --project-directory=. -f ../../go/docker-compose.yml -f ./docker-compose.override.yml

.PHONY: run-suite
run-suite: build down run

.PHONY: run
run:
	$(COMPOSE) up -d --no-build server
	make tmp

.PHONY: tmp
tmp:
	$(COMPOSE) run --no-deps wait-for-it -t 60 server:8000 -- echo "API is up"
	$(COMPOSE) run --no-deps webtests npm run docker-integration # first run
	$(COMPOSE) run --no-deps webtests npm run docker-integration # second run
	$(COMPOSE) run --no-deps webtests npm run docker-integration # third run

.PHONY: up
up:
	$(COMPOSE) up -d

.PHONY: down
down:
	$(COMPOSE) down --remove-orphans --volumes

.PHONY: logs
logs:
	$(COMPOSE) logs -f --tail=100

.PHONY: logs-once
logs-once:
	$(COMPOSE) logs

.PHONY: build-server
build-server:
	$(COMPOSE) build server

.PHONY: build-web
build-web:
	$(COMPOSE) build webtests

.PHONY: ps pull build
ps pull build:
	$(COMPOSE) $@
