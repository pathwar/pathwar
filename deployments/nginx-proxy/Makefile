.PHONY: all
all: up logs

.PHONY: up
up:
	-docker network create service-proxy
	docker-compose up -d
	-docker network connect service-proxy nginx-proxy

.PHONY: logs
logs:
	docker-compose logs --tail=100 -f
