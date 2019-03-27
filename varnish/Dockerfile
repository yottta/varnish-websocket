#FROM debian:jessie
# use the next one if debian has issues with the repositories
FROM ubuntu:16.04

RUN \
  useradd -r -s /bin/false varnishd

# Install Varnish source build dependencies.
RUN \
  apt-get update && apt-get install -y --no-install-recommends \
    automake \
    build-essential \
    ca-certificates \
    wget \
    python3-sphinx \
    libedit-dev \
    libjemalloc-dev \
    libncurses-dev \
    libpcre3-dev \
    libtool \
    pkg-config \
    python-docutils \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Install Varnish from source, so that Varnish modules can be compiled and installed.
ENV VARNISH_VERSION=6.2.0
RUN \
  apt-get update && \
  mkdir -p /usr/local/src && \
  cd /usr/local/src && \
  wget https://github.com/varnishcache/varnish-cache/archive/varnish-${VARNISH_VERSION}.tar.gz && \
  tar -xzf varnish-$VARNISH_VERSION.tar.gz && \
  cd varnish-cache-varnish-$VARNISH_VERSION && \
  ./autogen.sh && \
  ./configure && \
  make install && \
  rm ../varnish-$VARNISH_VERSION.tar.gz

COPY user.vcl /etc/varnish/user.vcl

WORKDIR /etc/varnish
ENTRYPOINT ["varnishd", "-F", "-a", ":6081", "-T", "localhost:6082", "-f", "/etc/varnish/user.vcl", "-s", "malloc,256m", "-p", "pipe_timeout=600"]
