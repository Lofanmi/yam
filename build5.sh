#!/bin/bash
ROOT=$(pwd)
SRC="${ROOT}/internal/php-src"

# php-src
if [ ! -d "${SRC}/5.6.40" ]; then
  if [ -d "./php-5.6.40" ]; then
    rm -rf ./php-5.6.40
  fi
  tar xvf "${SRC}/php-src-php-5.6.40.tar.xz"
  mv ./php-5.6.40 "${SRC}/5.6.40"
fi

# configure
if [ ! -f "${SRC}/5.6.40/Zend/zend_config.h" ]; then
  cd "${SRC}/5.6.40" && \
  ./configure --with-curl && \
  make -j8
fi

cd "$ROOT" && go build -tags=php5 -buildmode=c-shared -o ./bin/yamgo.so ./cmd/php5/
