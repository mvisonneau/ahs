##
# BUILD CONTAINER
##

FROM goreleaser/goreleaser:v0.165.0 as builder

WORKDIR /build

COPY . .
RUN \
apk add --no-cache make ;\
make build-linux-amd64

##
# RELEASE CONTAINER
##

FROM busybox:1.32-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/dist/ahs_linux_amd64/ahs /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/ahs"]
CMD [""]
