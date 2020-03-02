FROM    php:7.3-apache
RUN apt-get update -y && apt-get install -y libpng-dev
RUN docker-php-ext-install gd
COPY    www/ /var/www/html/
COPY    on-init /pwinit/
