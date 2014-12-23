# pathwar/nginx:1.7.8-onbuild pathwar/nginx:1.7-onbuild pathwar/nginx:1-onbuild pathwar/nginx:onbuild
FROM nginx:1.7.8
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
RUN wget -qO - http://j.mp/pathwar-scripts-install | bash

# Override pathwar scripts
RUN echo nginx -g daemon off > /pathwar/scripts/run

# -onbuild specific
ONBUILD CMD /pathwar/scripts/init
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD COPY www /usr/share/nginx/html/
