FROM golang:1.20-alpine as buildbase

RUN apk add build-base git

ARG CI_JOB_TOKEN

WORKDIR /go/src/github.com/rarimo/dashboard-rarime-link-svc

COPY . .

ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS="linux"

RUN go mod tidy
RUN go mod vendor
RUN go build -o /usr/local/bin/dashboard-rarime-link-svc github.com/rarimo/dashboard-rarime-link-svc

###

FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/dashboard-rarime-link-svc /usr/local/bin/dashboard-rarime-link-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["dashboard-rarime-link-svc"]
