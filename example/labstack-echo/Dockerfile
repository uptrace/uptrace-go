FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/labstack-echo

FROM base AS echo-server
RUN go install ./server.go
CMD ["/go/bin/server"]
