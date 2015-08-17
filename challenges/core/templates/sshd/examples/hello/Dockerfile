FROM pathwar/sshd:onbuild
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)
RUN apt-get update && apt-get install -y -q gcc && apt-get clean
COPY hello.c /home/bobby/
WORKDIR /home/bobby
RUN cc hello.c
RUN chown bobby:bobby /home/bobby/*
