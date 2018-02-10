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
