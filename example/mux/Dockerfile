FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/mux

FROM base AS mux-server
RUN go install ./server.go
CMD ["/go/bin/server"]
