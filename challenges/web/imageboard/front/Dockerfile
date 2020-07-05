FROM    php:7.3-apache
RUN     docker-php-ext-install mysqli
COPY    www/ /var/www/html/
COPY    on-init /pwinit/
