# pathwar/nginx:1.7.8-onbuild pathwar/nginx:1.7-onbuild pathwar/nginx:1-onbuild pathwar/nginx:onbuild
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

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /usr/share/nginx/html/
