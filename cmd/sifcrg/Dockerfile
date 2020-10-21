#
# Build
#
FROM golang:1.15 AS build

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

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

# Copy the compiled binairy over.
COPY --from=build /go/bin/sifcrg /usr/bin/sifcrg
