NAME =	pathwar/nginx
VERSIONS =	1.7.8 1.7.8-onbuild

all: build

example:	build-1.7.8-onbuild
	#$(MAKE) -C example build
	#$(MAKE) -C example run
	$(MAKE) -C example up

build:	build-1.7.8 build-1.7.8-onbuild

*/skeleton:	../../skeleton ../../skeleton/* ../../skeleton/*/*
	../../skeleton/install_local.sh $@

build-1.7.8:	1.7.8/skeleton
	docker build -t $(NAME):1.7.8 1.7.8
	docker tag $(NAME):1.7.8 $(NAME):1.7
	docker tag $(NAME):1.7.8 $(NAME):1
	docker tag $(NAME):1.7.8 $(NAME):latest

build-1.7.8-onbuild:	1.7.8-onbuild/skeleton
	docker build -t $(NAME):1.7.8-onbuild 1.7.8-onbuild
	docker tag $(NAME):1.7.8-onbuild $(NAME):1.7-onbuild
	docker tag $(NAME):1.7.8-onbuild $(NAME):onbuild
	docker tag $(NAME):1.7.8-onbuild $(NAME):1-onbuild

release:
	docker push $(NAME)
