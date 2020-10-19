#
# Build
#
FROM golang:1.15 AS build

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV SIF_CLI=sifnodecli
ENV SIF_DAEMON=sifnoded

# Empty dir for the db data
RUN mkdir /data

WORKDIR /sif
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make install

#
# Main
#
ENV PACKAGES curl jq bind-tools
FROM alpine

RUN apk add --update --no-cache $PACKAGES

# Copy the compiled binaires over.
COPY --from=build /go/bin/sifnoded /usr/bin/sifnoded
COPY --from=build /go/bin/sifnodecli /usr/bin/sifnodecli
COPY --from=build /go/bin/sifgen /usr/bin/sifgen

