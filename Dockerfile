ARG PHP_VERSION=7.4

FROM ghcr.io/aperim/php-apache-redirect:${PHP_VERSION}

ARG MAUTIC_VERSION=features
ARG COMPOSER_VERSION=1.10.20
ARG MAX_UPLOAD=256M
ARG MAX_MEMORY=512M
ARG MAX_TIME=180

LABEL org.opencontainers.image.source https://github.com/aperim/docker-mautic
LABEL org.label-schema.build-date=$BUILD_DATE \
  org.label-schema.name="Mautic v${MAUTIC_VERSION} in a container" \
  org.label-schema.description="Mautic version ${MAUTIC_VERSION} in a container" \
  org.label-schema.url="https://github.com/aperim/docker-mautic" \
  org.label-schema.vcs-ref=$VCS_REF \
  org.label-schema.vcs-url="https://github.com/aperim/docker-mautic" \
  org.label-schema.vendor="Aperim" \
  org.label-schema.version=$VERSION \
  org.label-schema.schema-version="1.0"

COPY ./docker-entrypoint.sh /usr/sbin/docker-entrypoint.sh

RUN apt-get update && apt-get install -y \
  curl \
  git \
  libzip-dev \
  libpng-dev \
  libwebp-dev \
  libjpeg-dev \
  libc-client-dev \
  libkrb5-dev \
  unzip \
  tzdata \
  cron && \
  rm -rf /var/lib/apt/lists/* && \
  mkdir -p /usr/src/composer /etc/cron.d /usr/src/mautic/${MAUTIC_VERSION} && \
  curl -sS https://getcomposer.org/installer -o /usr/src/composer/composer-setup.php && \
  cd /usr/src/composer/ && \
  php composer-setup.php --install-dir=/usr/local/bin --filename=composer && \
  composer self-update && \
  docker-php-ext-install pdo zip bcmath sockets gd mysqli pdo_mysql && \
  docker-php-ext-configure imap --with-kerberos --with-imap-ssl && docker-php-ext-install imap && \
  composer self-update ${COMPOSER_VERSION} && \
  chmod +x /usr/sbin/docker-entrypoint.sh && \
  printf "memory_limit = ${MAX_MEMORY}\n\nupload_max_filesize = ${MAX_UPLOAD}\npost_max_size = ${MAX_UPLOAD} \nmax_file_uploads = ${MAX_UPLOAD}\n\nmax_input_time = ${MAX_TIME}\nmax_execution_time = ${MAX_TIME}\n\nmax_input_vars = 5000" > /usr/local/etc/php/conf.d/mautic.ini && \
  touch /etc/cron.d/mautic && \
  mkfifo /var/log/cron.pipe && \
  chown -R www-data:www-data /usr/src/mautic /usr/sbin/docker-entrypoint.sh /etc/cron.d/mautic /var/log/cron.pipe 

RUN git clone --single-branch --branch ${MAUTIC_VERSION} https://github.com/mautic/mautic.git /var/www/html && \
    chown -R www-data:www-data /var/www/html && \
    cd /var/www/html && \
    su -s /bin/bash www-data -c "composer install"

RUN cd /var/www/html && \
    su -s /bin/bash www-data -c "tar zcvf /usr/src/mautic/${MAUTIC_VERSION}/www.tgz ." && \
    curl -L https://raw.githubusercontent.com/mautic/docker-mautic/master/apache/makeconfig.php -o /usr/src/mautic/${MAUTIC_VERSION}/makeconfig.php

USER www-data:www-data
ENV MAUTIC_VERSION=${MAUTIC_VERSION}
EXPOSE 80
WORKDIR /var/www/html
ENTRYPOINT [ "/usr/sbin/docker-entrypoint.sh" ]
CMD [ "/usr/sbin/apache2",  "-DFOREGROUND"]