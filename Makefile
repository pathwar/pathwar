## moul/rules.mk
DOCKER_IMAGE = pathwar/pathwar
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

.PHONY: integration
integration:
	cd tool/integration; make run-suite


##
## docker.build
##

# rules.mk will already build the original Dockerfile
docker.build: agent.docker.build

api.docker.push: api.docker.build
api.docker.push:
	docker push pathwar/pathwar
agent.docker.push: agent.docker.build
agent.docker.push:
	docker push pathwar/agent

.PHONY: agent.docker.build
agent.docker.build:
	$(call docker_build,Dockerfile.agent,pathwar/agent)
.PHONY: api.docker.build
api.docker.build:
	$(call docker_build,Dockerfile,pathwar/pathwar)

.PHONY: agent.deploy
agent.deploy: agent.docker.push
	ssh $(PATHWAR_AGENT_DEPLOY_HOST) 'cd $(PATHWAR_AGENT_DEPLOY_PATH) && make down pull up logs'
.PHONY: api.deploy
api.deploy: api.docker.push
	ssh $(PATHWAR_API_DEPLOY_HOST) 'cd $(PATHWAR_API_DEPLOY_PATH) && make pull up-fast logs'
