FROM    php:7.3-apache
RUN     apt-get update \
        && apt-get install -y libicu-dev \
        && docker-php-ext-configure intl \
        && docker-php-ext-install intl \
        && docker-php-ext-install gettext

RUN     apt-get install -y locales-all
COPY    php.ini /usr/local/etc/php/conf.d/php.ini
COPY    translations /var/www/translations/
COPY    www/ /var/www/html/
COPY    on-init /pwinit/
