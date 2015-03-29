NAME =	hello

all:	run

build:
	docker-compose build

shell:	build
	docker-compose run $(NAME) /bin/bash

up:	build
	docker-compose up
