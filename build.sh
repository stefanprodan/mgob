#! /bin/sh

apk add --no-cache ca-certificates tzdata bash curl krb5-dev

# Install GnuPG
if [ "_${MGOB_EN_GPG}" = "_true" ]
then
  apk add gnupg=${GNUPG_VERSION}
fi

cd /tmp

# Install MinIO
if [ "_${MGOB_EN_MINIO}" = "_true" ]
then
  curl -O https://dl.minio.io/client/mc/release/linux-amd64/mc
  mv mc /usr/bin
  chmod u+x /usr/bin/mc
fi

# Install RClone
if [ "_${MGOB_EN_RCLONE}" = "_true" ]
then
  curl -O https://downloads.rclone.org/rclone-current-linux-amd64.zip
  unzip rclone-current-linux-amd64.zip
  cp rclone-*-linux-amd64/rclone /usr/bin/
  chmod u+x /usr/bin/rclone
  rm rclone-current-linux-amd64.zip
fi

#install gcloud
if [ "_${MGOB_EN_GCLOUD}" = "_true" ]
then
  export PATH="/google-cloud-sdk/bin:$PATH"
  apk --no-cache add \
        python3 \
        py3-pip \
        libc6-compat \
        openssh-client \
        git
  pip3 --no-cache-dir install --upgrade pip
  pip --no-cache-dir install wheel
  pip --no-cache-dir install crcmod
  curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
  tar xzf google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
  mv google-cloud-sdk /
  rm google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
  ln -s /lib /lib64
  gcloud config set core/disable_usage_reporting true
  gcloud config set component_manager/disable_update_check true
  gcloud config set metrics/environment github_docker_image
  gcloud --version
fi

# install azure-cli and aws-cli
if [ "_${MGOB_EN_AZURE}" = "_true" -o "_${MGOB_EN_AWS_CLI}" = "_true" ]
then
  apk --no-cache add python3 py3-pip
  apk --no-cache add --virtual=build gcc libffi-dev musl-dev openssl-dev  python3-dev make
  pip3 --no-cache-dir install --upgrade pip
  pip --no-cache-dir install wheel cffi
  echo "EN_AZURE: $MGOB_EN_AZURE; EN_AWS_CLI: $MGOB_EN_AWS_CLI"
  [ "_${MGOB_EN_AZURE}" = "_true" ] && pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}
  [ "_${MGOB_EN_AWS_CLI}" = "_true" ] && pip --no-cache-dir install awscli==${AWS_CLI_VERSION}
  apk del --purge build
fi
