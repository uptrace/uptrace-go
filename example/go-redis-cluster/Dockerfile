FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/go-redis-cluster

FROM base AS go-redis-cluster-client
RUN go install ./client.go
CMD ["/go/bin/client"]
