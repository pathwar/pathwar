.PHONY: all
all: up logs

.PHONY: ps restart pull
ps restart pull:
	docker-compose $@

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down --remove-orphans

.PHONY: logs
logs:
	docker-compose logs -f --tail=1000

.PHONY: flushdb
flushdb: down
	docker volume rm platform_db_data
