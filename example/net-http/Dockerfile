FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/net-http

FROM base AS http-server
RUN go install ./server.go
CMD ["/go/bin/server"]
