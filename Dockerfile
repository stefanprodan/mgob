ARG MONGODB_TOOLS_VERSION=100.5.1
ARG EN_AWS_CLI=true
ARG AWS_CLI_VERSION=1.22.46
ARG EN_AZURE=true
ARG AZURE_CLI_VERSION=2.32.0
ARG EN_GCLOUD=true
ARG GOOGLE_CLOUD_SDK_VERSION=370.0.0
ARG EN_GPG=true
ARG GNUPG_VERSION="2.2.31-r0"
ARG EN_MINIO=true
ARG EN_RCLONE=true

FROM golang:1.17 as mgob-builder

ARG VERSION

COPY . /go/src/github.com/stefanprodan/mgob

WORKDIR /go/src/github.com/stefanprodan/mgob

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags "-X main.version=$VERSION" \
    -a -installsuffix cgo \
    -o mgob github.com/stefanprodan/mgob/cmd/mgob

FROM golang:1.17-alpine3.15 as tools-builder

ARG MONGODB_TOOLS_VERSION

RUN apk add --no-cache git build-base krb5-dev \
    && git clone https://github.com/mongodb/mongo-tools.git --depth 1 --branch $MONGODB_TOOLS_VERSION

WORKDIR mongo-tools
RUN ./make build

FROM alpine:3.15

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
ARG MONGODB_TOOLS_VERSION
ARG AWS_CLI_VERSION
ARG AZURE_CLI_VERSION
ARG GOOGLE_CLOUD_SDK_VERSION
ARG GNUPG_VERSION
ARG EN_AWS_CLI
ARG EN_AZURE
ARG EN_GCLOUD
ARG EN_GPG
ARG EN_MINIO
ARG EN_RCLONE

ENV MONGODB_TOOLS_VERSION=$MONGODB_TOOLS_VERSION \
    GNUPG_VERSION=$GNUPG_VERSION \
    GOOGLE_CLOUD_SDK_VERSION=$GOOGLE_CLOUD_SDK_VERSION \
    AZURE_CLI_VERSION=$AZURE_CLI_VERSION \
    AWS_CLI_VERSION=$AWS_CLI_VERSION \
    MGOB_EN_AWS_CLI=$EN_AWS_CLI \
    MGOB_EN_AZURE=$EN_AZURE \
    MGOB_EN_GCLOUD=$EN_GCLOUD \
    MGOB_EN_GPG=$EN_GPG \
    MGOB_EN_MINIO=$EN_MINIO \
    MGOB_EN_RCLONE=$EN_RCLONE

WORKDIR /

COPY build.sh /tmp
RUN /tmp/build.sh

COPY --from=mgob-builder /go/src/github.com/stefanprodan/mgob/mgob .
COPY --from=tools-builder /go/mongo-tools/bin/* /usr/bin/

VOLUME ["/config", "/storage", "/tmp", "/data"]

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="mgob" \
      org.label-schema.description="MongoDB backup automation tool" \
      org.label-schema.url="https://github.com/stefanprodan/mgob" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/stefanprodan/mgob" \
      org.label-schema.vendor="stefanprodan.com" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"

ENTRYPOINT [ "./mgob" ]
