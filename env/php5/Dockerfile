
FROM php:5.6-fpm

LABEL maintainer="Lofanmi <Lofanmi@gmail.com>"

ENV TZ=Asia/Shanghai

RUN sed -i "s@http://deb.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list && \
    sed -i "s@http://security.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list && \
    apt update -y && \
    apt install -y cmake git libcurl4-openssl-dev libxml2-dev wget xz-utils zip

ADD https://mirrors.aliyun.com/composer/composer.phar /usr/bin/composer

RUN pecl channel-update pecl.php.net && \
    chmod +x /usr/bin/composer && \
    composer config -g repo.packagist composer 'https://mirrors.aliyun.com/composer/'

ENV GOPATH='/root/go'
ENV GOROOT='/usr/local/go'
ENV GOPROXY='https://goproxy.cn,direct'
ENV GOSUMDB='sum.golang.google.cn'
ENV GO111MODULE='on'
ENV PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

RUN wget -O /usr/local/go1.17.7.linux-arm64.tar.gz -c https://studygolang.com/dl/golang/go1.17.7.linux-arm64.tar.gz && \
    cd /usr/local && \
    tar xvf go1.17.7.linux-arm64.tar.gz && \
    rm -rf /usr/local/go1.17.7.linux-arm64.tar.gz && \
    cd /tmp && mkdir -p $GOPATH/bin

WORKDIR /app

CMD [ "bash" ]