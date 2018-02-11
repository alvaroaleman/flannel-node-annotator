FROM alpine:3.7

COPY flannel-node-annotator /usr/local/bin

USER nobody

CMD ["/usr/local/bin/flannel-node-annotator", "-logtostderr"]
