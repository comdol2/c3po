FROM hub.docker.io/golang:1.13-alpine as builder
	
ARG VER=1.0

ENV GITHUB github.com
ENV GOPATH /go
ENV GOBIN /workspace
ENV GOSRC $GOPATH/src
ENV GOGIT $GOSRC/$GITHUB/comdol2/snow

RUN mkdir -p $GOGIT $GOBIN

WORKDIR $GOGIT

COPY . $GOGIT/

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags "-s -w -X $GITHUB/comdol2/snow/cmd.version=${VER}" -a -o $GOBIN/snow

##################################################

FROM hub.docker.io/alpine:3.11

WORKDIR /

COPY --from=builder /workspace/snow .

CMD ["/bin/sh"]
