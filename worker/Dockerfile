# Build stage
FROM golang:1.11 as builder

WORKDIR /usr/bin/

WORKDIR /go/src/search-worker
COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

# ldflags "-s -w" strips binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -installsuffix cgo -o search-worker


# Release stage
FROM alpine:3.8

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/search-worker/search-worker   .

ENV PATH=$PATH:/root/

CMD ["search-worker"]