FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/echo

FROM base AS echo-server
RUN go install ./server.go
CMD ["/go/bin/server"]
