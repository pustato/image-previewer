FROM golang:1.17.9-alpine

RUN mkdir -p /opt/integration_tests
WORKDIR /opt/integration_tests

RUN adduser -u 1000 -s /bin/bash -S -D intgrtest
USER intgrtest

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
CMD ["go", "test", "-tags=intgrtest", "./..."]
