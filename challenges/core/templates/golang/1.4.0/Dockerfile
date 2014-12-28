# pathwar/golang:1.4.0 pathwar/golang:1.4 pathwar/golang:1 pathwar/golang:latest
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
