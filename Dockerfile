FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git
ENV CGO_ENABLED=0
RUN go install github.com/minio/minio@latest && \
    go install github.com/minio/mc@latest

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/minio /usr/local/bin/minio
COPY --from=builder /go/bin/mc /usr/local/bin/mc

EXPOSE 9000 9001
VOLUME ["/data"]

CMD ["minio", "server", "/data", "--console-address", ":9001", "--address", ":9000"]