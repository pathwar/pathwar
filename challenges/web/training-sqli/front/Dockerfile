FROM    php:7.3-apache
RUN     apt-get update \
        && docker-php-ext-install mysqli \
        && apt-get install -y libicu-dev \
        && docker-php-ext-configure intl \
        && docker-php-ext-install intl \
        && docker-php-ext-install gettext

#locales
RUN     apt-get install -y locales-all

COPY    translations/ /var/www/translations/
COPY    www/ /var/www/html/
COPY    on-init /pwinit/
