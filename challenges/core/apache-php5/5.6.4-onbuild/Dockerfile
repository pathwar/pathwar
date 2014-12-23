# pathwar/apache-php:5.6.4-onbuild pathwar/apache-php:5.6-onbuild pathwar/apache-php:5-onbuild pathwar/apache-php:onbuild
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
#RUN curl bit.ly/pathwar-scripts-install | VERSION=1.0 bash

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /var/www/html/
