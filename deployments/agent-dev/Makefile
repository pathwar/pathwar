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

.PHONY: auth
auth:
	docker-compose run --rm agent-dev --debug agent --no-run

.PHONY: clean
clean:
	docker-compose run --rm agent-dev --debug agent --no-run --clean

.PHONY: nginx-config
nginx-config:
	docker exec pathwar-agent-nginx cat /etc/nginx/nginx.conf

.PHONY: fclean
fclean:
	docker kill `docker ps -q`
	docker system prune -f

.PHONY: agent-ps
agent-ps:
	docker-compose run --rm agent-dev --debug composez ps
