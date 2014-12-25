# pathwar/apache-php:5.6.4-mysqli-onbuild pathwar/apache-php:5.6-mysqli-onbuild pathwar/apache-php:5-mysqli-onbuild pathwar/apache-php:mysqli-onbuild
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl mysql-client \
 && apt-get clean

# PHP dependencies
RUN docker-php-ext-install mysqli

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /var/www/html/
