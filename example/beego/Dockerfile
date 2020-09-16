FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/beego

FROM base AS beego-server
RUN go install ./server.go
CMD ["/go/bin/server"]
