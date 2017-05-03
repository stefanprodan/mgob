FROM alpine:edge

RUN apk add --no-cache mongodb-tools ca-certificates
ADD https://dl.minio.io/client/mc/release/linux-amd64/mc /usr/bin
RUN chmod u+x /usr/bin/mc

WORKDIR /root/
COPY mgob    .

VOLUME ["config", "storage", "tmp"]

ENTRYPOINT [ "./mgob" ]