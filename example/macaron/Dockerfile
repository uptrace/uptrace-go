FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/macaron

FROM base AS macaron-server
RUN go install ./server.go
CMD ["/go/bin/server"]
