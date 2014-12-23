# pathwar/apache-php:5.6.4 pathwar/apache-php:5.6 pathwar/apache-php:5 pathwar/apache-php:latest
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl \
 && apt-get clean

# Add Pathwar tools
RUN curl -Lk http://j.mp/pathwar-scripts-install | DL=curl bash

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run
