# pathwar/golang:1.4.0-onbuild pathwar/golang:1.4-onbuild pathwar/golang:1-onbuild pathwar/golang:onbuild
FROM golang:1.4.0
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
#RUN apt-get update \
# && apt-get -y -q upgrade \
# && apt-get -y -q install curl \
# && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo go-wrapper run > /pathwar/scripts/run

# Onbuild
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
ONBUILD COPY . /go/src/app
ONBUILD RUN go-wrapper download
ONBUILD RUN go-wrapper install
