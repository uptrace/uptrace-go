FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/go-restful

FROM base AS restful-server
RUN go install ./server.go
CMD ["/go/bin/server"]
