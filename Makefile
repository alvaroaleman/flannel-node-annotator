SHELL = /bin/bash

REGISTRY ?= docker.io
REGISTRY_NAMESPACE ?= alvaroaleman

IMAGE_TAG ?= latest
IMAGE_NAME = $(REGISTRY)/$(REGISTRY_NAMESPACE)/flannel-node-annotator:$(IMAGE_TAG)

vendor: Gopkg.lock Gopkg.toml
	dep ensure -vendor-only

node-annotator: $(shell find controller -name '*.go') main.go vendor
		@docker run --rm \
			-v $$PWD:/go/src/github.com/alvaroaleman/k8s-node-pulicip-annotator \
			-w /go/src/github.com/alvaroaleman/k8s-node-pulicip-annotator \
			golang:1.9.4 \
			env CGO_ENABLED=0 go build \
				-ldflags '-s -w' \
				-o node-annotator \
				github.com/alvaroaleman/k8s-node-pulicip-annotator

docker-image:
	docker build -t $(IMAGE_NAME) .
	docker push $(IMAGE_NAME)
