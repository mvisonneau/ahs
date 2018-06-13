##
# BUILD CONTAINER
##

FROM golang:1.10 as builder

WORKDIR /go/src/github.com/mvisonneau/ahs

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make deps ;\
make build

##
# RELEASE CONTAINER
##

FROM scratch

WORKDIR /

COPY --from=builder /go/src/github.com/mvisonneau/ahs /

ENTRYPOINT ["/ahs"]
CMD [""]
