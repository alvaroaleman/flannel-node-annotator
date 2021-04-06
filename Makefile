SHELL = /bin/bash

REGISTRY ?= docker.io
REGISTRY_NAMESPACE ?= alvaroaleman

IMAGE_TAG ?= latest
IMAGE_NAME = $(REGISTRY)/$(REGISTRY_NAMESPACE)/flannel-node-annotator:$(IMAGE_TAG)

gotidy:
	go mod tidy

flannel-node-annotator: $(shell find controller -name '*.go') main.go go.mod go.sum
		@docker run --rm \
			-v $$PWD:/go/src/github.com/alvaroaleman/flannel-node-annotator \
			-w /go/src/github.com/alvaroaleman/flannel-node-annotator \
			golang:1.16.2 \
			env CGO_ENABLED=0 go build \
				-ldflags '-s -w' \
				-o flannel-node-annotator \
				github.com/alvaroaleman/flannel-node-annotator

docker-image: flannel-node-annotator
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)

.PHONY: docker-image gotidy
