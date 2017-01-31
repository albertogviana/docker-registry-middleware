FROM alpine
MAINTAINER Alberto Guimaraes Viana <albertogviana@gmail.com>

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    apk --update add ca-certificates

EXPOSE 8080
CMD ["docker-registry-middleware"]

COPY docker-registry-middleware /usr/local/bin/docker-registry-middleware
RUN chmod +x /usr/local/bin/docker-registry-middleware

ARG user=admin
ARG group=admin
ARG uid=1000
ARG gid=1000

# Admin is run with user `admin`, uid = 1000
# If you bind mount a volume from the host or a data container,
# ensure you use the same uid
RUN addgroup -g ${gid} ${group} && \
    adduser -h "/home/admin" -u ${uid} -G ${group} -s /bin/sh -D ${user}

USER admin
