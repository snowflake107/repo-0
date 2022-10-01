ARG PHP_VERSION=8.0
FROM php:${PHP_VERSION}-apache

ARG USER=www-data
ARG GROUP=www-data
ARG COMPOSER_SHA=55ce33d7678c5a611085589f1f3ddf8b3c52d662cd01d4ba75c0ee0459970c2200a51f492d557530c71c15d8dba01eae
ARG MAUTIC_VERSION=4.4.3

RUN apt-get update && \
    apt-get -y full-upgrade && \
    apt-get install --no-install-recommends -y \
        cron=3.0pl1-137 \
        git=1:2.30.2-1 \
        libc-client2007e-dev=8:2007f~dfsg-7+b1 \
        libcurl4-openssl-dev=7.74.0-1.3+deb11u3 \
        libc-client2007e-dev=8:2007f~dfsg-7+b1 \
        libicu-dev=67.1-7 \
        libkrb5-dev=1.18.3-6+deb11u2 \
        libmcrypt-dev=2.5.8-3.4+b1 \
        libonig-dev=6.9.6-1.1 \
        libpng-dev=1.6.37-3 \
        zlib1g-dev=1:1.2.11.dfsg-2+deb11u2 \
        libssl-dev=1.1.1n-0+deb11u3 \
        libfreetype6-dev=2.10.4+dfsg-1+deb11u1 \
        libjpeg-dev=1:2.0.6-4 \
        libzip-dev=1.7.3-1 \
        sudo=1.9.5p2-3 \
        unzip=6.0-26+deb11u1 \
        zip=3.0-12 && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -r /var/lib/apt/lists/* && \
    rm /etc/cron.daily/* && \
    chmod gu+rw /var/run && \
    chmod gu+s /usr/sbin/cron && \
    docker-php-ext-configure imap --with-kerberos --with-imap-ssl && \
    docker-php-ext-configure gd --with-freetype --with-jpeg && \
    docker-php-ext-configure opcache --enable-opcache && \
    docker-php-ext-install bcmath gd imap intl mbstring mysqli curl pdo_mysql zip opcache bcmath gd sockets exif && \
    docker-php-ext-enable imap intl mbstring mysqli curl pdo_mysql zip sodium opcache bcmath gd sockets exif && \
    a2dismod mpm_event && \
    a2enmod mpm_prefork rewrite && \
    mkdir -p /opt/composer /var/www/.composer /usr/src/mautic-cron /usr/src/mautic-php /var/sessions/php /var/run && \
    cd /opt/composer && \
    php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"  && \
    php -r "if (hash_file('sha384', 'composer-setup.php') === '${COMPOSER_SHA}') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;"  && \
    php composer-setup.php && \
    php -r "unlink('composer-setup.php');" && \
    mv composer.phar /usr/local/bin/composer && \
    git clone --branch ${MAUTIC_VERSION} --depth 1 https://github.com/mautic/mautic.git /usr/src/mautic && \
    rm -Rf /usr/src/mautic/.git && \
    chown -R ${USER}:${GROUP} /usr/src/mautic /var/www/.composer && \
    mkfifo /var/log/cron.pipe && \
    cp /usr/local/etc/php/php.ini-production /usr/local/etc/php/php.ini && \
    cp /usr/local/etc/php/php.ini-production /usr/src/mautic-php/php.ini && \
    echo "pcre.jit = 0" > /usr/local/etc/php/conf.d/mautic-php-disable-pcre-jit.ini && \
    echo "always_populate_raw_post_data = -1" > /usr/local/etc/php/conf.d/mautic-php-always-populate-raw-post.ini && \
    touch /var/run/crond.pid && \
    chown www-data:www-data /var/log/cron.pipe \
        /usr/local/etc/php/php.ini \
        /usr/src/mautic-php/php.ini \
        /usr/local/etc/php/conf.d/mautic-php-disable-pcre-jit.ini \
        /usr/local/etc/php/conf.d/mautic-php-always-populate-raw-post.ini \
        /var/sessions/php \
        /var/run/crond.pid && \
    chmod 444 /usr/local/etc/php/conf.d/mautic-php-disable-pcre-jit.ini \
        /usr/local/etc/php/conf.d/mautic-php-always-populate-raw-post.ini \
        /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^\;date\.timezone =.*$/date.timezone = "\${PHP_INI_DATE_TIMEZONE}"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^memory_limit = .*$/memory_limit = "\${PHP_MEMORY_LIMIT}"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^upload_max_filesize = .*$/upload_max_filesize = "\${PHP_MAX_UPLOAD}"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^post_max_size = .*$/post_max_size = "\${PHP_MAX_UPLOAD}"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^max_execution_time = .*$/max_execution_time = "\${PHP_MAX_EXECUTION_TIME}"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^variables_order = .*$/variables_order = "EGPCS"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^;session\.save_path = .*$/session.save_path = "\/var\/sessions\/php"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^;error_log = php_errors\.log$/error_log = "\/dev\/stderr"/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^error_reporting = .*$/error_reporting = E_ALL/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^display_errors = .*$/display_errors = On/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^max_execution_time = .*$/max_execution_time = 900/g' /usr/src/mautic-php/php.ini && \
    sed --in-place --regexp-extended 's/^<VirtualHost \*\:80>$/<VirtualHost \*\:\${PORT}>/g' /etc/apache2/sites-available/000-default.conf && \
    sed --in-place --regexp-extended 's/^Listen 80$/Listen \${PORT}/g' /etc/apache2/ports.conf

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN curl -L https://raw.githubusercontent.com/mautic/docker-mautic/mautic4/common/docker-entrypoint.sh -o /docker-entrypoint.sh && \
    curl -L https://raw.githubusercontent.com/mautic/docker-mautic/mautic4/common/makeconfig.php -o /makeconfig.php && \
    curl -L https://raw.githubusercontent.com/mautic/docker-mautic/mautic4/common/makedb.php -o /makedb.php && \
    curl -L https://raw.githubusercontent.com/mautic/docker-mautic/mautic4/common/mautic.crontab -o /usr/src/mautic-cron/mautic && \
    sed --in-place --regexp-extended 's/^(    sudo -Eu www-data php \/var\/www\/html\/bin\/console mautic:install)/    echo >\&2 \"Params: \${INSTALL_PARAMS[*]}\"\n\1/g' /docker-entrypoint.sh && \
    sed --in-place --regexp-extended 's/^(php \/makeconfig.php)/\1\ncat >\&2 \/var\/www\/html\/app\/config\/local.php/g' /docker-entrypoint.sh && \
    sed --in-place --regexp-extended 's/^(\$path\s+=\s+\x27\/var\/www\/html\/app\/config\/local\.php\x27\;)/\1\nfwrite(\$stderr, json_encode\(\$_ENV, JSON_PRETTY_PRINT\)\)\;\n/g' /makeconfig.php && \
    cp /usr/src/mautic-cron/mautic /etc/cron.d/mautic && \
    chmod 755 /docker-entrypoint.sh /makeconfig.php /makedb.php && \
    rm -f /usr/src/mautic/sig.txt && \
    find /usr/src/mautic -type f | sort | tr '\n' '\0' | xargs -0 sha1sum | sha1sum | head -c 40 > /usr/src/mautic/sig.txt && \
    chown www-data:www-data /usr/src/mautic/sig.txt && \
    chmod 400 /usr/src/mautic/sig.txt && \
    echo "Signature for build: $(cat /usr/src/mautic/sig.txt)" && \
    chown www-data:www-data /etc/cron.d/mautic /usr/src/mautic-cron/mautic && \
    chmod 644 /etc/cron.d/mautic /usr/src/mautic-cron/mautic

COPY rootfs /

USER ${USER}

RUN cd /usr/src/mautic && \
    composer install && \
    cp /usr/src/mautic-php/php.ini /usr/local/etc/php/php.ini

HEALTHCHECK --interval=30s \
    --timeout=30s \
    --start-period=2h \
    --retries=3 \
    CMD docker-mautic-healthcheck

ARG BUILD_DESCRIPTION="Mautic version ${MAUTIC_VERSION} that works in a container"
ARG BUILD_NAME="Mautic ${MAUTIC_VERSION}"
ARG BUILD_DATE
ARG BUILD_REF
ARG BUILD_REPOSITORY
ARG BUILD_VERSION

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

ARG MAUTIC_RUN_CRON_JOBS=true \
    MAUTIC_RUN_MIGRATIONS=false \
    MAUTIC_DB_USER=root \
    MAUTIC_DB_HOST=localhost \
    MAUTIC_DB_NAME=mautic \
    MAUTIC_DB_PORT=3306 \
    MAUTIC_HOME=/var/www/html \
    PHP_INI_DATE_TIMEZONE=UTC \
    PHP_MEMORY_LIMIT=512M \
    PHP_MAX_UPLOAD=512M \
    PHP_MAX_EXECUTION_TIME=300 \
    PORT=8080

ENV MAUTIC_RUN_CRON_JOBS=${MAUTIC_RUN_CRON_JOBS} \
    MAUTIC_RUN_MIGRATIONS=${MAUTIC_RUN_MIGRATIONS} \
    MAUTIC_DB_HOST=${MAUTIC_DB_HOST} \
    MAUTIC_DB_USER=${MAUTIC_DB_USER} \
    MAUTIC_DB_NAME=${MAUTIC_DB_NAME} \
    MAUTIC_DB_PORT=${MAUTIC_DB_PORT} \
    MAUTIC_HOME=${MAUTIC_HOME} \
    PHP_INI_DATE_TIMEZONE=${PHP_INI_DATE_TIMEZONE} \
    PHP_MEMORY_LIMIT=${PHP_MEMORY_LIMIT} \
    PHP_MAX_UPLOAD=${PHP_MAX_UPLOAD} \
    PHP_MAX_EXECUTION_TIME=${PHP_MAX_EXECUTION_TIME} \
    PORT=${PORT}

STOPSIGNAL SIGWINCH
WORKDIR ${MAUTIC_HOME}
EXPOSE ${PORT}
ENTRYPOINT ["docker-mautic-entrypoint"]
CMD ["apache2-foreground"]
