# pathwar/nginx:1.7.8 pathwar/nginx:1.7 pathwar/nginx:1 pathwar/nginx:latest
FROM nginx:1.7.8
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
RUN curl bit.ly/pathwar-scripts-install | bash

# Update pathwar scripts
RUN echo nginx -g daemon off > /pathwar/scripts/run
