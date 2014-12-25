# pathwar/apache-php:full-5.6.4 pathwar/apache-php:full-5.6 pathwar/apache-php:full-5 pathwar/apache-php:full
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl mysql-client libpng12-dev \
 && apt-get clean

# PHP dependencies
RUN docker-php-ext-install gd mysql mysqli

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run
