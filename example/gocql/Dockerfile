FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/gocql

FROM base AS gocql-client
RUN go install ./client.go
CMD ["/go/bin/client"]
