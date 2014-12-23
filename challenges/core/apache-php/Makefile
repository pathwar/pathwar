NAME =	pathwar/apache-php

build:
	docker build -t $(NAME):5.6.4 5.6.4
	docker build -t $(NAME):5.6.4-onbuild 5.6.4-onbuild

release:
	docker tag $(NAME):5.6.4-onbuild $(NAME):5.6-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):5-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):onbuild
	docker tag $(NAME):5.6.4 $(NAME):5.6
	docker tag $(NAME):5.6.4 $(NAME):5
	docker tag $(NAME):5.6.4 $(NAME):latest
	docker push $(NAME)
