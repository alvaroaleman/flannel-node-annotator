SHELL = /bin/bash

REGISTRY ?= docker.io
REGISTRY_NAMESPACE ?= alvaroaleman

IMAGE_TAG ?= latest
IMAGE_NAME = $(REGISTRY)/$(REGISTRY_NAMESPACE)/flannel-node-annotator:$(IMAGE_TAG)

DOCKER ?= docker

vendor: Gopkg.lock Gopkg.toml
	dep ensure -vendor-only

flannel-node-annotator: $(shell find controller -name '*.go') main.go vendor
		@$(DOCKER) run --rm \
			-v $$PWD:/go/src/github.com/alvaroaleman/flannel-node-annotator \
			-w /go/src/github.com/alvaroaleman/flannel-node-annotator \
			golang:1.9.4 \
			env CGO_ENABLED=0 go build \
				-ldflags '-s -w' \
				-o flannel-node-annotator \
				github.com/alvaroaleman/flannel-node-annotator

docker-image:
	$(DOCKER) build -t $(IMAGE_NAME) .
	$(DOCKER) push $(IMAGE_NAME)
