FROM php:7-apache

ARG BUILD_DESCRIPTION="Mautic that works in a container"
ARG BUILD_NAME="Mautic"
ARG BUILD_DATE
ARG BUILD_REF
ARG BUILD_REPOSITORY
ARG BUILD_VERSION

RUN apt-get update && \
    apt-get -y full-upgrade && \
    apt-get -y install \
        git=1-2.30.2-1 \
        libc-client-dev=8-2007f~dfsg-7+b1 \
        libkrb5-dev=1.18.3-6+deb11u1 \
        libpng-dev=1.6.37-3 \
        zlib1g-dev-1=1.2.11.dfsg-2 \
        libzip-dev=1.7.3-1 && \
    rm -r /var/lib/apt/lists/* && \
    docker-php-ext-configure imap --with-kerberos --with-imap-ssl && \
    docker-php-ext-install bcmath gd imap zip sockets && \
    mkdir -p /opt/composer && \
    cd /opt/composer && \
    php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"  && \
    php -r "if (hash_file('sha384', 'composer-setup.php') === '906a84df04cea2aa72f40b5f787e49f22d4c2f19492ac310e8cba5b96ac8b64115ac402c8cd292b8a03482574915d1a8') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;"  && \
    php composer-setup.php && \
    php -r "unlink('composer-setup.php');" && \
    mv composer.phar /usr/local/bin/composer && \
    git clone https://github.com/mautic/mautic.git /app && \
    cd /app && \
    composer install

COPY rootfs /

HEALTHCHECK --interval=30s \
    --timeout=30s \
    --start-period=4m \
    --retries=3 \
    CMD curl -f http://localhost:80

LABEL \
    maintainer="Troy Kelly <troy@aperim.com>" \
    org.opencontainers.image.title="${BUILD_NAME}" \
    org.opencontainers.image.description="${BUILD_DESCRIPTION}" \
    org.opencontainers.image.vendor="Troy Kelly" \
    org.opencontainers.image.authors="Troy Kelly <troy@aperim.com>" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.url="https://aperim.com" \
    org.opencontainers.image.source="https://github.com/${BUILD_REPOSITORY}" \
    org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/README.md" \
    org.opencontainers.image.created=${BUILD_DATE} \
    org.opencontainers.image.revision=${BUILD_REF} \
    org.opencontainers.image.version=${BUILD_VERSION}

ENTRYPOINT ["docker-mautic-entrypoint"]