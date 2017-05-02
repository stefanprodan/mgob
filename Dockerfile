FROM alpine:edge

RUN apk add --no-cache mongodb-tools

WORKDIR /root/
COPY mgob    .

VOLUME ["config", "storage", "tmp"]

ENTRYPOINT [ "./mgob" ]