FROM    php:7.3-apache
RUN apt-get update -y && apt-get install -y libpng-dev
RUN docker-php-ext-install gd
RUN apt-get install -y libicu-dev \
        && docker-php-ext-configure intl \
        && docker-php-ext-install intl \
        && docker-php-ext-install gettext \
        && apt-get install -y locales-all

COPY    translations /var/www/translations/ 
COPY    www/ /var/www/html/
COPY    on-init /pwinit/
