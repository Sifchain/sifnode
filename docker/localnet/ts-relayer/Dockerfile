#
# Build
#
FROM golang:1.15 AS build

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /sif
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make install

#
# Main
#
FROM alpine

ENV PACKAGES curl jq bind-tools expect nodejs npm

RUN apk add --update --no-cache $PACKAGES

COPY --from=build /go/bin/sifnoded /usr/bin/sifnoded

RUN npm i -g @confio/relayer@main
