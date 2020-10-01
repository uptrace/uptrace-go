FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/grpc

FROM base AS grpc-client
RUN go install ./client/client.go
CMD ["/go/bin/client"]
