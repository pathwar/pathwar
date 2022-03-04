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
	bash -ec 'source .env; set -x; docker-compose run --rm agent-dev --debug agent --no-run --salt=$$PATHWAR_AGENT_SALT --moderator-password=$$PATHWAR_AGENT_MODERATOR_PASSWORD --domain-suffix=$$PATHWAR_AGENT_NGINX_DOMAIN_SUFFIX --agent-name=$$PATHWAR_AGENT_NAME --port=$$PATHWAR_AGENT_NGINX_PORT --sso-realm=$$PATHWAR_AGENT_SSO_REALM --sso-token-file=./agent-token.json --sso-pubkey="$$PATHWAR_AGENT_SSO_PUBKEY" --http-api-addr=$$PATHWAR_AGENT_HTTP_API_ADDR'


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
