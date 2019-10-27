.PHONY: all
all: up logs

.PHONY: up
up:
	docker-compose up -d

.PHONY: logs
logs:
	docker-compose logs -f --tail=1000
