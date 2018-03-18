FROM ubuntu:latest
LABEL maintainer="r@countableset.com"

RUN apt update
RUN apt install -y imagemagick libtiff-tools sane-utils
COPY document-imaging /opt/
RUN mkdir /scan

VOLUME /scan
WORKDIR /scan

CMD [ "sh", "-c", "/opt/document-imaging" ]
