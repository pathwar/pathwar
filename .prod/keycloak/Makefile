.PHONY: all
all: up logs

.PHONY: up
up:
	docker-compose up -d

.PHONY: logs
logs:
	docker-compose logs -f --tail=1000

.PHONY: mysql-shell
mysql-shell:
	. ./.env; docker-compose exec keycloakdb mysql -p$$MYSQL_ROOT_PASSWORD keycloak
