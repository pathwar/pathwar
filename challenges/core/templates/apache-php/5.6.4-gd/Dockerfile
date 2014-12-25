# pathwar/apache-php:5.6.4-gd pathwar/apache-php:5.6-gd pathwar/apache-php:5-gd pathwar/apache-php:gd
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
