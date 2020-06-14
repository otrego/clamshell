# This uses the distroless images as the base:
# https://github.com/GoogleContainerTools/distroless

FROM golang:1.14 as build-env

WORKDIR /go/src/github.com/otrego/clamshell
ADD . /go/src/github.com/otrego/clamshell

RUN go get -d -v ./...

RUN go build -o /go/bin/github.com/clamshell/server /go/src/github.com/otrego/clamshell/server

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/github.com/clamshell/server /
CMD ["/server"]
