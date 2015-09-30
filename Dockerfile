FROM python:3.5-slim

RUN apt-get update && apt-get install -y cmake pkg-config
ADD https://github.com/libgit2/libgit2/archive/v0.23.2.tar.gz /
RUN cd / && tar xzf v0.23.2.tar.gz && cd libgit2-0.23.2/ && cmake . && make && make install

RUN ldconfig

ADD https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz /usr/local/
RUN cd /usr/local && tar zxf go1.5.1.linux-amd64.tar.gz && rm go1.5.1.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
ENV GO15VENDOREXPERIMENT 1
