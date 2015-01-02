NAME =	server

build:
	fig build

up:	build
	fig up -d
	fig logs

shell:	build
	fig run $(NAME) /bin/bash
