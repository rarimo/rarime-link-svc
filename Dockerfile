FROM golang:1.22-alpine as buildbase

RUN apk add build-base git

ARG CI_JOB_TOKEN

WORKDIR /go/src/github.com/rarimo/rarime-link-svc

COPY . .

ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS="linux"

RUN go mod tidy
RUN go mod vendor
RUN go build -o /usr/local/bin/rarime-link-svc github.com/rarimo/rarime-link-svc

###

FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/rarime-link-svc /usr/local/bin/rarime-link-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["rarime-link-svc"]
