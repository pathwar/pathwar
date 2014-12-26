# pathwar/node:0.10.34-onbuild pathwar/node:0.10-onbuild pathwar/node:0-onbuild pathwar/node:onbuild
FROM node:0.10.34
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
#RUN apt-get update \
# && apt-get -y -q upgrade \
# && apt-get -y -q install curl \
# && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo npm start > /pathwar/scripts/run

# Onbuild
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ONBUILD COPY app/package.json /usr/src/app/
ONBUILD RUN npm install
ONBUILD COPY app /usr/src/app
