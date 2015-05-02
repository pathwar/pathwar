FROM pathwar/supervisord:onbuild
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)
EXPOSE 22 80
RUN apt-get update \
 && apt-get install -y -q openssh-server \
 && apt-get clean
