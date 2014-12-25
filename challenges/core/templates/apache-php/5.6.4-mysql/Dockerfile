# pathwar/apache-php:5.6.4-mysql pathwar/apache-php:5.6-mysql pathwar/apache-php:5-mysql pathwar/apache-php:mysql
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl mysql-client \
 && apt-get clean

# PHP dependencies
RUN docker-php-ext-install mysql

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run
