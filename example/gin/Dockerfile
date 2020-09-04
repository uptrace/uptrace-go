FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/gin

FROM base AS gin-server
RUN go install ./server.go
CMD ["/go/bin/server"]
