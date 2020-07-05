docker.build pathwar.prepare pathwar.push:
	for dir in $(dir $(wildcard ./*/*/docker-compose.yml)); do set -e && cd $$dir && make $@ && cd -; done

_ci: docker.build
