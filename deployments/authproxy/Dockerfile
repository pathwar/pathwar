FROM pathwar/nginx-lua


MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)


RUN apt-get update \
 && apt-get upgrade -q -y \
 && apt-get install -q -y \
        lua-cjson \
        lua-socket \
        lua-sec \
        python-pip
RUN pip install requests


COPY lua /pathwar/lua
COPY conf_generation/ /pathwar/conf_generation/
COPY nginx.conf /etc/nginx/nginx.conf


RUN chmod +x /pathwar/conf_generation/generate_conf.py


CMD python /pathwar/conf_generation/generate_conf.py && /usr/sbin/nginx -g 'daemon off;'


ENV PATHWAR_API_SCHEME https://
ENV PATHWAR_API_HOST api.pathwar.net
ENV PATHWAR_API_USER root
