NAME =	pathwar/nginx

build:
	docker build -t $(NAME):1.7.8 1.7.8
	docker build -t $(NAME):1.7.8-onbuild 1.7.8-onbuild

release:
	docker tag $(NAME):1.7.8-onbuild $(NAME):1.7-onbuild
	docker tag $(NAME):1.7.8-onbuild $(NAME):onbuild
	docker tag $(NAME):1.7.8-onbuild $(NAME):1-onbuild
	docker tag $(NAME):1.7.8 $(NAME):1.7
	docker tag $(NAME):1.7.8 $(NAME):1
	docker tag $(NAME):1.7.8 $(NAME):latest
	docker push $(NAME)
