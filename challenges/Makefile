CHALLENGES := $(dir $(wildcard ./*/*/docker-compose.yml))

docker.build pathwar.prepare pathwar.push:
	for dir in $(CHALLENGES); do set -e; \
	  cd $$dir && make $@ && cd -; \
	done

_ci: docker.build
