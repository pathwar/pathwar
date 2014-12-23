# pathwar/apache-php:5.6.4 pathwar/apache-php:5.6 pathwar/apache-php:5 pathwar/apache-php:latest
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
RUN wget -qO - http://j.mp/pathwar-scripts-install | bash

# Override pathwar scripts
RUN apache2 -DFOREGROUND > /pathwar/scripts/run
