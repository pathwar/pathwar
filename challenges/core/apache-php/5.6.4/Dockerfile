# pathwar/apache-php:5.6.4 pathwar/apache-php:5.6 pathwar/apache-php:5 pathwar/apache-php:latest
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
RUN curl bit.ly/pathwar-scripts-install | VERSION=1.0 bash

# Override pathwar scripts
RUN apache2 -DFOREGROUND > /pathwar/scripts/run
