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
RUN node > /pathwar/scripts/run
