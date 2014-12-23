NAME =	pathwar/apache-php

build:
	docker build -t $(NAME):5.6 5.6
	docker build -t $(NAME):5.6-onbuild 5.6-onbuild

release:
	docker tag $(NAME):5.6-onbuild $(NAME):5-onbuild
	docker tag $(NAME):5.6-onbuild $(NAME):onbuild
	docker tag $(NAME):5.6 $(NAME):5
	docker tag $(NAME):5.6 $(NAME):latest
	docker push $(NAME)
