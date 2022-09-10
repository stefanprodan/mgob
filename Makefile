SHELL:=/bin/bash

APP_VERSION?=1.5

# build vars
BUILD_DATE:=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")


#run vars
CONFIG:=$$(pwd)/test/config
TRAVIS:=$$(pwd)/test/travis

# go tools
PACKAGES:=$(shell go list ./... | grep -v '/vendor/')
VETARGS:=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -rangeloops -shift -structtags -unsafeptr

build:
	@echo ">>> Building $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)"
	@docker build \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg VCS_REF=$(GITHUB_SHA) \
	    --build-arg VERSION=$(APP_VERSION) \
	    -t $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION) .

aws:
	@echo ">>> Building $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)"
	@docker build \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg VCS_REF=$(GITHUB_SHA) \
	    --build-arg VERSION=$(APP_VERSION) \
	    --build-arg EN_AWS_CLI=true \
	    --build-arg EN_AZURE=false \
	    --build-arg EN_GCLOUD=false \
	    --build-arg EN_MINIO=false \
	    --build-arg EN_RCLONE=false \
	    --build-arg EN_GPG=true \
	    -t $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)-aws .

azure:
	@echo ">>> Building $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)"
	@docker build \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg VCS_REF=$(GITHUB_SHA) \
	    --build-arg VERSION=$(APP_VERSION) \
	    --build-arg EN_AWS_CLI=false \
	    --build-arg EN_AZURE=true \
	    --build-arg EN_GCLOUD=false \
	    --build-arg EN_MINIO=false \
	    --build-arg EN_RCLONE=false \
	    --build-arg EN_GPG=true \
	    -t $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)-azure .

gcloud:
	@echo ">>> Building $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)"
	@docker build \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg VCS_REF=$(GITHUB_SHA) \
	    --build-arg VERSION=$(APP_VERSION) \
	    --build-arg EN_AWS_CLI=false \
	    --build-arg EN_AZURE=false \
	    --build-arg EN_GCLOUD=true \
	    --build-arg EN_MINIO=false \
	    --build-arg EN_RCLONE=false \
	    --build-arg EN_GPG=true \
	    -t $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)-gcloud .

ci:
	@echo ">>> Building mgob:$(APP_VERSION).$(GITHUB_RUN_NUMBER) image"
	@docker build \
	    --build-arg BUILD_DATE=$(BUILD_DATE) \
	    --build-arg VCS_REF=$(GITHUB_SHA) \
	    --build-arg VERSION=$(APP_VERSION).$(GITHUB_RUN_NUMBER) \
	    -t $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION).$(GITHUB_RUN_NUMBER) .

	@echo ">>> Starting mgob container"
	@docker run -d --net=host --name mgob \
	    --restart unless-stopped \
	    -v "$(TRAVIS):/config" \
	    -v "/tmp/ssh_host_rsa_key:/etc/ssh/ssh_host_rsa_key:ro" \
	    -v "/tmp/ssh_host_rsa_key.pub:/etc/ssh/ssh_host_rsa_key.pub:ro" \
        $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION).$(GITHUB_RUN_NUMBER) \
		-ConfigPath=/config \
		-StoragePath=/storage \
		-TmpPath=/tmp \
		-LogLevel=info

publish:
	@echo $(DOCKER_PASS) | docker login -u "$(DOCKER_USER)" --password-stdin
	@docker tag $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION).$(GITHUB_RUN_NUMBER) $(GITHUB_REPOSITORY_OWNER)/mgob:edge
	@docker push $(GITHUB_REPOSITORY_OWNER)/mgob:edge

release:
	@echo $(DOCKER_PASS) | docker login -u "$(DOCKER_USER)" --password-stdin
	@docker tag $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION).$(GITHUB_RUN_NUMBER) $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)
	@docker tag $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION).$(GITHUB_RUN_NUMBER) $(GITHUB_REPOSITORY_OWNER)/mgob:latest
	@docker push $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION)
	@docker push $(GITHUB_REPOSITORY_OWNER)/mgob:latest

run:
	@echo ">>> Starting mgob container"
	@docker run -dp 8090:8090 --name mgob-$(APP_VERSION) \
	    --restart unless-stopped \
	    -v "$(CONFIG):/config" \
        $(GITHUB_REPOSITORY_OWNER)/mgob:$(APP_VERSION) \
		-ConfigPath=/config \
		-StoragePath=/storage \
		-TmpPath=/tmp \
		-LogLevel=info

backend:
	@docker run -dp 20022:22 --name mgob-sftp \
	    atmoz/sftp:alpine test:test:::backup
	@docker run -dp 20099:9000 --name mgob-s3 \
	    -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE" \
	    -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
	    minio/minio server /export
	@mc config host add local http://localhost:20099 \
	    AKIAIOSFODNN7EXAMPLE wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY S3v4
	@sleep 5
	@mc mb local/backup

fmt:
	@echo ">>> Running go fmt $(PACKAGES)"
	@go fmt $(PACKAGES)

vet:
	@echo ">>> Running go vet $(VETARGS)"
	@go list ./... \
		| grep -v /vendor/ \
		| cut -d '/' -f 4- \
		| xargs -n1 \
			go tool vet $(VETARGS) ;\
	if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "go vet failed"; \
	fi

.PHONY: build
