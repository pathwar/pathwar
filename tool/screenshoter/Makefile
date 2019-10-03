URL ?= https://www.pathwar.land/

.PHONY: run
run:
	docker-compose up -d chrome-debugger
	docker-compose run screenshoter -screenshot-url='$(URL)'
