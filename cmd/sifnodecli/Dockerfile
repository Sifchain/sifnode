#
# Build
#
FROM golang:1.15 AS build

ARG chainnet

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV CHAINNET=$chainnet

# Empty dir for the db data
RUN mkdir /data

WORKDIR /sif
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CHAINNET=$CHAINNET make install

#
# Main
#
FROM alpine

ENV PACKAGES curl jq bind-tools

RUN apk add --update --no-cache $PACKAGES

# Copy the compiled binairy over.
COPY --from=build /go/bin/sifnodecli /usr/bin/sifnodecli
