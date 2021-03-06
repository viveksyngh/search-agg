# Build stage
FROM golang:1.11 as builder

WORKDIR /usr/bin/

WORKDIR /go/src/github.com/viveksyngh/search-agg
COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

# ldflags "-s -w" strips binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -installsuffix cgo -o search-api


# Release stage
FROM alpine:3.8

RUN apk --no-cache add ca-certificates

EXPOSE 8000

WORKDIR /root/

COPY --from=builder /go/src/github.com/viveksyngh/search-agg/search-api   .

ENV PATH=$PATH:/root/

CMD ["search-api"]