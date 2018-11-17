FROM snapcore/snapcraft:latest

# Install tools and libraries
RUN apt-get update && \
    apt-get install -y \
    imagemagick \
    libtiff-tools \
    sane-utils \
    devscripts \
    dh-make \
    wget

# Install Golang
ENV GOLANG_VERSION 1.11.2
RUN wget https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

VOLUME /document-imaging
WORKDIR /document-imaging