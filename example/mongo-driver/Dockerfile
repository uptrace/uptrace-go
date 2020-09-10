FROM golang:alpine AS base
COPY . /src/
WORKDIR /src/example/mongo-driver

FROM base AS mongo-driver
RUN go install ./mongo.go
CMD ["/go/bin/mongo"]
