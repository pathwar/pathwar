.PHONY: build
build:
	@for dir in `find . -type f -name docker-compose.yml | sed 's@/[^/]*$$@@' | sort | uniq`; do (set -xe; \
	  cd $$dir; \
	  docker-compose build --pull; \
	); done
