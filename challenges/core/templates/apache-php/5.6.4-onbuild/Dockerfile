# pathwar/apache-php:5.6.4-onbuild pathwar/apache-php:5.6-onbuild pathwar/apache-php:5-onbuild pathwar/apache-php:onbuild
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl mysql-client libpng12-dev \
 && apt-get clean

# PHP dependencies
RUN docker-php-ext-install gd mysql mysqli \
 && rm -f /usr/local/etc/php/conf.d/*

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /var/www/html/
ONBUILD COPY php.ini /usr/local/etc/php/conf.d/
#ONBUILD RUN chown -R www-data:www-data /var/www/html/; chmod -R 777 /var/www/html/
