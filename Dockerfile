FROM alpine:3.7

COPY node-annotator /usr/local/bin

USER nobody

CMD ["/usr/local/bin/node-annotator", "-logtostderr"]
