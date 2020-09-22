FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/gomemcache

FROM base AS gomemcache-client
RUN go install ./client.go
CMD ["/go/bin/client"]
