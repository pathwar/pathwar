# pathwar/apache-php:5.6.4-mysqli pathwar/apache-php:5.6-mysqli pathwar/apache-php:5-mysqli pathwar/apache-php:mysqli
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
