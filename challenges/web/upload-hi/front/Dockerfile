FROM    php:7.3-apache
#RUN apt-get update && docker-php-ext-install \
#    mbstring \
#    gettext \
#    intl
RUN apt-get -y update \
    && apt-get install -y libicu-dev \
    && docker-php-ext-configure intl \
    && docker-php-ext-install intl \
    && docker-php-ext-install gettext


COPY    translations /var/www/translations/
COPY    www/ /var/www/html/
COPY    on-init /pwinit/
