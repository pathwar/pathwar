# pathwar/supervisord:latest
FROM ubuntu:wily
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
RUN apt-get update \
 && apt-get -y -q upgrade \
 && apt-get -y -q install curl supervisor \
 && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo /usr/bin/supervisord > /pathwar/scripts/run
