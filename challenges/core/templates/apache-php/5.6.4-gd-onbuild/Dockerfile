# pathwar/apache-php:5.6.4-gd-onbuild pathwar/apache-php:5.6-gd-onbuild pathwar/apache-php:5-gd-onbuild pathwar/apache-php:gd-onbuild
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl libpng12-dev \
 && apt-get clean

# PHP dependencies
RUN docker-php-ext-install gd

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /var/www/html/
