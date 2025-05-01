FROM golang:1.24.2-alpine3.21 AS builder
WORKDIR /go/src/int-bluesky
COPY . .
RUN \
    apk add protoc protobuf-dev make git && \
    make build

FROM alpine:3.21.3
RUN apk --no-cache add ca-certificates \
    && update-ca-certificates
COPY --from=builder /go/src/int-bluesky/int-bluesky /bin/int-bluesky
ENTRYPOINT ["/bin/int-bluesky"]
