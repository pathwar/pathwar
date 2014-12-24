NAME =	pathwar/apache-php
VERSIONS =	5.6.4 5.6.4-onbuild

all: build

example:	build-5.6.4-onbuild
	#$(MAKE) -C example build
	#$(MAKE) -C example run
	$(MAKE) -C example up

build:	build-5.6.4 build-5.6.4-onbuild

*/skeleton:	../../skeleton ../../skeleton/* ../../skeleton/*/*
	../../skeleton/install_local.sh $@

build-5.6.4:	5.6.4/skeleton
	docker build -t $(NAME):5.6.4 5.6.4
	docker tag $(NAME):5.6.4 $(NAME):5.6
	docker tag $(NAME):5.6.4 $(NAME):5
	docker tag $(NAME):5.6.4 $(NAME):latest

build-5.6.4-onbuild:	5.6.4-onbuild/skeleton
	docker build -t $(NAME):5.6.4-onbuild 5.6.4-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):5.6-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):5-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):onbuild

release:
	docker push $(NAME)
