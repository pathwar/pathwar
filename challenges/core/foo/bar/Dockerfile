FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# RUN wget | bash ...
CMD /pathwar/scripts/init

COPY scripts /pathwar/scripts/
COPY level.yml /pathwar/

COPY www /var/www/html/
