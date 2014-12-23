NAME =	pathwar/apache-php

build:
	docker build -t $(NAME):tmp .

release:	build
	docker tag $(NAME):tmp $(NAME):5.6-onbuild
	docker tag $(NAME):tmp $(NAME):5-onbuild
	docker tag $(NAME):tmp $(NAME):onbuild
	docker rmi $(NAME):tmp
	docker push $(NAME)
