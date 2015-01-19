all:	build

build:
	fig build

run:
	fig stop
	fig up -d
	fig ps
	fig logs
