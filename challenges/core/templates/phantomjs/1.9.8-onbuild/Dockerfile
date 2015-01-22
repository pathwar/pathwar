# pathwar/phantomjs:1.9.8-onbuild
# pathwar/phantomjs:1.9-onbuild
# pathwar/phantomjs:1-onbuild
# pathwar/phantomjs:onbuild

FROM rdpanek/phantomjs:1.9.8
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo /usr/bin/phantomjs main.js > /pathwar/scripts/run

# Onbuild
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD CMD /pathwar/scripts/init

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ONBUILD COPY app /usr/src/app
