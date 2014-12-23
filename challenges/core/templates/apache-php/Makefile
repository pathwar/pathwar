NAME =	pathwar/apache-php

build:
	../../skeleton/install_local.sh 5.6.4/skeleton
	docker build -t $(NAME):5.6.4 5.6.4
	docker tag $(NAME):5.6.4 $(NAME):5.6
	docker tag $(NAME):5.6.4 $(NAME):5
	docker tag $(NAME):5.6.4 $(NAME):latest
	../../skeleton/install_local.sh 5.6.4-onbuild/skeleton
	docker build -t $(NAME):5.6.4-onbuild 5.6.4-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):5.6-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):5-onbuild
	docker tag $(NAME):5.6.4-onbuild $(NAME):onbuild

release:
	docker push $(NAME)
