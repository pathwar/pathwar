## moul/rules.mk
INSTALL_STEPS += recursive_install
GENERATE_STEPS += recursive_generate
TEST_STEPS += recursive_test
#BUMPDEPS_STEPS += recursive_bumpdeps
include rules.mk  # see https://github.com/moul/rules.mk
##

.PHONY: recursive_install
recursive_install:
	cd go; make install

.PHONY: recursive_test
recursive_test:
	cd go; make test
	cd web; make test

.PHONY: recursive_bumpdeps
recursive_bumpdeps:
	cd go; make bumpdeps

.PHONY: clean
clean:
	cd go; make clean
	cd docs; make clean

.PHONY: recursive_generate
recursive_generate:
	cd go; make generate
	cd docs; make generate


##
## Integration
##

## FIXME: TODO
#    .PHONY: integration
#    integration: integration.build integration.run
#
#    .PHONY: integration.build
#    integration.build:
#    	docker-compose build server web
#
#    .PHONY:integration.run
#    integration.run:
#    	docker-compose up -d --no-build server
#    	docker-compose exec server ./wait-for-it.sh serverdb:3306 -- echo serverdb ready
#    	docker-compose exec server ./wait-for-it.sh localhost:9111 -- echo gRPC ready
#    	sleep 5
#    	#docker-compose exec server pathwar.land sql adduser --sql-config=$$SQL_CONFIG --email=integration@example.com --username=integration#     --password=integration
#    	docker-compose run web npm test
#
#    .PHONY: generate-fake-data
#    generate-fake-data:
#    	AUTH_TOKEN=`http --check-status :8000/authenticate username=integration | jq -r .token` && \
#    	  http POST :8000/dev/generate-fake-data Authorization:$$AUTH_TOKEN && \
#    	  http POST :8000/dev/sql-dump Authorization:$$AUTH_TOKEN
