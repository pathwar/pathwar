FROM pathwar/unix:onbuild
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)
COPY level.c /home/bobby/level.c
RUN touch .passwd                                     \
 && cc -o level level.c                               \
 && chown super-bobby:bobby ./level ./level.c .       \
 && chown super-bobby:super-bobby .passwd             \
 && chmod 600 .passwd                                 \
 && chmod 750 ./level .                               \
 && chmod u+s level
