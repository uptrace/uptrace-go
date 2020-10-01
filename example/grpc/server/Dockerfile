FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/grpc

FROM base AS grpc-server
RUN go install ./server/server.go
CMD ["/go/bin/server"]
