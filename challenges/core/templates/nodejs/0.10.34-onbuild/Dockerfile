# pathwar/apache-php:5.6.4 pathwar/apache-php:5.6 pathwar/apache-php:5 pathwar/apache-php:latest
FROM node:0.10.34
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl \
 && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo npm start > /pathwar/scripts/run

# Onbuild
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY app/package.json /usr/src/app/
ONBUILD RUN npm install
ONBUILD COPY app /usr/src/app
