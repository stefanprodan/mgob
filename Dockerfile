FROM golang:1.16-buster

ARG VERSION

COPY . /go/src/github.com/stefanprodan/mgob

WORKDIR /go/src/github.com/stefanprodan/mgob

RUN CGO_ENABLED=0 GOOS=linux \
      go build \
        -ldflags "-X main.version=$VERSION" \
        -a -installsuffix cgo \
        -o mgob github.com/stefanprodan/mgob/cmd/mgob

FROM ubuntu:focal

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION

ENV MONGODB_TOOLS_VERSION 100.5.1
ENV MONGODB_VERSION 4.4
ENV GNUPG_VERSION 2.2.31-r0
ENV GOOGLE_CLOUD_SDK_VERSION 364.0.0
ENV AZURE_CLI_VERSION 2.30.0
ENV AWS_CLI_VERSION 2.4.0
ENV PATH /root/google-cloud-sdk/bin:$PATH

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="mgob" \
      org.label-schema.description="MongoDB backup automation tool" \
      org.label-schema.url="https://github.com/stefanprodan/mgob" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/stefanprodan/mgob" \
      org.label-schema.vendor="stefanprodan.com" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"

RUN apt-get update && apt install apt-transport-https ca-certificates gnupg tzdata wget unzip curl -y
RUN wget -qO - https://www.mongodb.org/static/pgp/server-${MONGODB_VERSION}.asc | apt-key add -
RUN echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/${MONGODB_VERSION} multiverse" | tee /etc/apt/sources.list.d/mongodb-org-${MONGODB_VERSION}.list
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

RUN apt update && apt install -y google-cloud-sdk
RUN gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image && \
    gcloud --version

RUN wget https://fastdl.mongodb.org/tools/db/mongodb-database-tools-ubuntu2004-x86_64-100.5.1.tgz
RUN tar -zxvf mongodb-database-tools-*-100.5.1.tgz
RUN cp /mongodb-database-tools-*-100.5.1/bin/* /usr/bin/

ADD https://dl.minio.io/client/mc/release/linux-amd64/mc /usr/bin
RUN chmod u+x /usr/bin/mc

ADD https://downloads.rclone.org/rclone-current-linux-amd64.zip /tmp
RUN cd /tmp \
  && unzip rclone-current-linux-amd64.zip \
  && cp rclone-*-linux-amd64/rclone /usr/bin/ \
  && chmod u+x /usr/bin/rclone

WORKDIR /root/

# install azure-cli and aws-cli
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${AWS_CLI_VERSION}.zip" -o "awscliv2.zip" && unzip awscliv2.zip && ./aws/install

RUN apt autoremove -y
COPY --from=0 /go/src/github.com/stefanprodan/mgob/mgob .

VOLUME ["/config", "/storage", "/tmp", "/data"]

ENTRYPOINT [ "./mgob" ]
