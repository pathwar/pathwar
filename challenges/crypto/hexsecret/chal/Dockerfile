FROM    nginx:1.22-alpine
RUN     apk add --no-cache python3
COPY    default.conf.template /etc/nginx/templates/
RUN     mkdir -p /chal
COPY    on-init /pwinit/
