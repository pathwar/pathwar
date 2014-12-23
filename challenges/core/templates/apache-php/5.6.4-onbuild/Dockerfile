# pathwar/apache-php:5.6.4-onbuild pathwar/apache-php:5.6-onbuild pathwar/apache-php:5-onbuild pathwar/apache-php:onbuild
FROM php:5.6-apache
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
#RUN apt-get update \
# && apt-get -y -q upgrade \
# && apt-get -y -q install curl \
# && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar
#RUN curl -Lk http://j.mp/pathwar-scripts-install | DL=curl bash

# Override pathwar scripts
RUN echo apache2 -DFOREGROUND > /pathwar/scripts/run

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /var/www/html/
