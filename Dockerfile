##
# BUILD CONTAINER
##

FROM goreleaser/goreleaser:v0.133.0 as builder

WORKDIR /build

COPY . .
RUN \
apk add --no-cache make ;\
make build-linux-amd64

##
# RELEASE CONTAINER
##

FROM busybox:1.31-glibc

WORKDIR /

COPY --from=builder /build/dist/ahs_linux_amd64/ahs /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/ahs"]
CMD [""]
