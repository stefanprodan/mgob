SHELL:=/bin/bash

APP_VERSION?=0.0.1

# build vars
BUILD_DATE:=$(shell date -u +%Y-%m-%d_%H.%M.%S)
GIT_REPOSITORY:=github.com/stefanprodan/mgob
GIT_COMMIT:=$(shell git rev-parse HEAD)
GIT_BRANCH:=$(shell git symbolic-ref --short HEAD)
MAINTAINER:="Stefan Prodan"
REPOSITORY?=stefanprodan

#run vars
CONFIG:=$$(pwd)/test/config
TRAVIS:=$$(pwd)/test/travis

# go tools
PACKAGES:=$(shell go list ./... | grep -v '/vendor/')
VETARGS:=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -rangeloops -shift -structtags -unsafeptr

travis:
	@echo ">>> Building mgob:build image"
	@docker build --build-arg APP_VERSION=$(APP_VERSION) -t $(REPOSITORY)/mgob:build -f Dockerfile.build .
	@docker create --name mgob_extract $(REPOSITORY)/mgob:build
	@docker cp mgob_extract:/dist/mgob ./mgob
	@docker rm -f mgob_extract
	@echo ">>> Building mgob:$(APP_VERSION) image"
	@docker build -t $(REPOSITORY)/mgob:$(APP_VERSION) .
	@rm ./mgob
	@echo ">>> Starting mgob container"
	@docker run -dp 8090:8090 --name mgob-$(APP_VERSION) \
	    --restart unless-stopped \
	    -v "$(TRAVIS):/config" \
        $(REPOSITORY)/mgob:$(APP_VERSION) \
		-ConfigPath=/config \
		-StoragePath=/storage \
		-TmpPath=/tmp \
		-LogLevel=info
	@curl http://localhost:8090/version

run: build
	@echo ">>> Starting mgob container"
	@docker run -dp 8090:8090 --name mgob-$(APP_VERSION) \
	    --restart unless-stopped \
	    -v "$(CONFIG):/config" \
        $(REPOSITORY)/mgob:$(APP_VERSION) \
		-ConfigPath=/config \
		-StoragePath=/storage \
		-TmpPath=/tmp \
		-LogLevel=info

build: clean
	@echo ">>> Building mgob:build image"
	@docker build --build-arg APP_VERSION=$(APP_VERSION) -t $(REPOSITORY)/mgob:build -f Dockerfile.build .
	@docker create --name mgob_extract $(REPOSITORY)/mgob:build
	@docker cp mgob_extract:/dist/mgob ./mgob
	@docker rm -f mgob_extract
	@echo ">>> Building mgob:$(APP_VERSION) image"
	@docker build -t $(REPOSITORY)/mgob:$(APP_VERSION) .
	@rm ./mgob

clean:
	@docker rm -f mgob-$(APP_VERSION) || true
	@docker rmi $$(docker images | awk '$$1 ~ /mgob/ { print $$3 }') || true
	@docker volume rm $$(docker volume ls -f dangling=true -q)

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

