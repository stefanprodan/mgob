FROM alpine:edge

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="mgob" \
      org.label-schema.description="MongoDB backup automation tool" \
      org.label-schema.url="https://github.com/stefanprodan/mgob" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/stefanprodan/mgob" \
      org.label-schema.vendor="stefanprodan.com" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"

RUN apk add --no-cache mongodb-tools ca-certificates
ADD https://dl.minio.io/client/mc/release/linux-amd64/mc /usr/bin
RUN chmod u+x /usr/bin/mc

WORKDIR /root/
COPY mgob    .

VOLUME ["/config", "/storage", "/tmp"]

ENTRYPOINT [ "./mgob" ]