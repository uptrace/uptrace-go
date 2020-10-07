FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/go-redis

FROM base AS go-redis-client
RUN go install ./client.go
CMD ["/go/bin/client"]
