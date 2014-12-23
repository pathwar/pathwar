# pathwar/nginx:1.7.8 pathwar/nginx:1.7 pathwar/nginx:1 pathwar/nginx:latest
FROM nginx:1.7.8
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Install dependencies
#RUN apt-get update \
# && apt-get -y -q upgrade \
# && apt-get -y -q install curl \
# && apt-get clean

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo 'nginx -g "daemon off;"' > /pathwar/scripts/run
