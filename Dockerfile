FROM golang:1.11-alpine3.8 AS build

RUN apk add --no-cache curl git build-base && \
	curl -SL -o /usr/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
	chmod +x /usr/bin/dep

ENV MAINDIR $GOPATH/src/github.com/go-ocf/certificate-authority
ENV APP certificate-authority-service
WORKDIR $MAINDIR
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -v --vendor-only
COPY . .

FROM build AS build-secure
RUN OCF_INSECURE=false go generate ./vendor/github.com/go-ocf/kit/security/
WORKDIR $MAINDIR/cmd/$APP
RUN go build -o /go/bin/$APP
WORKDIR $MAINDIR

FROM alpine:3.8 AS certificate-authority
COPY --from=build-secure /go/bin/$APP /
ENTRYPOINT ["/certificate-authority-service"]

FROM build AS build-insecure
RUN OCF_INSECURE=true go generate ./vendor/github.com/go-ocf/kit/security/
WORKDIR $MAINDIR/cmd/$APP
RUN go build -o /go/bin/$APP
WORKDIR $MAINDIR

FROM alpine:3.8 AS certificate-authority-insecure
COPY --from=build-insecure /go/bin/$APP /
ENTRYPOINT ["/certificate-authority-service"]
