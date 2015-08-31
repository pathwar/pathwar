FROM pathwar/unix:onbuild
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)
COPY level.c /home/bobby/
WORKDIR /home/bobby
RUN cc -o level level.c
RUN chown bobby:bobby /home/bobby/*
