FROM pathwar/supervisord:onbuild
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)
EXPOSE 22 80

RUN apt-get update \
 && apt-get install -y -q openssh-server \
 && apt-get clean

RUN mkdir -p /var/run/sshd /var/log/supervisor \
 && sed -i '/^PermitRootLogin/d;$s/$/\nPermitRootLogin yes/' /etc/ssh/sshd_config \
 && echo "root:root" | chpasswd \
 && chsh -s /bin/bash www-data \
 && echo "www-data:www-data" | chpasswd \
