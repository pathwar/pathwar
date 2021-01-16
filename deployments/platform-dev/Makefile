.PHONY: all
all: up logs

.PHONY: ps restart pull
ps restart pull:
	docker-compose $@

.PHONY: up-fast
up-fast:
	docker-compose up -d

.PHONY: up
up:
	docker-compose up -d mailserver pathwar-dev
	# wait for the migrations to finish
	sleep 120
	docker-compose up -d

.PHONY: down
down:
	docker-compose down --remove-orphans

.PHONY: logs
logs:
	docker-compose logs -f --tail=1000

.PHONY: flushdb
flushdb: down
	docker volume rm -f platform_db_data platform-dev_db_data || true

.PHONY: sqldb
sqldb:
	docker-compose exec db sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE'

.PHONY: sqldump
sqldump:
	docker-compose exec db sh -c 'mysqldump -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE'

.PHONY: backup
backup:
	@mkdir -p backups/`date +%Y-%m`
	docker-compose exec db sh -c 'mysqldump -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE --skip-extended-insert' > backups/`date +%Y-%m`/`date +%Y-%m-%d_%H:%M:%S`.mysql

BACKUP_PATH ?= $(shell ls -Art backups/*/*.mysql | tail -n 1)

.PHONY: restore
restore:
	docker-compose exec -T db sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE' < "$(BACKUP_PATH)"
