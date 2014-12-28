# pathwar/ruby:2.1.5-onbuild pathwar/ruby:2.1-onbuild pathwar/ruby:2-onbuild pathwar/ruby:onbuild
FROM ruby:2.1.5
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
#RUN apt-get update \
# && apt-get -y -q upgrade \
# && apt-get -y -q install curl \
# && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo irb > /pathwar/scripts/run

# Onbuild
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/

RUN bundle config --global frozen 1
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ONBUILD COPY app/Gemfile /usr/src/app/
ONBUILD COPY app/Gemfile.lock /usr/src/app/
ONBUILD RUN bundle install
ONBUILD COPY app /usr/src/app
