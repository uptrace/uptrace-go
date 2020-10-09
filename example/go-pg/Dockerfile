FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/go-pg

FROM base AS go-pg-client
RUN go install ./client.go
CMD ["/go/bin/client"]
