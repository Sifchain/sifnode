#
# Build
#
FROM golang:1.15.2 AS build

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=1
ENV GOOS=linux

# Empty dir for the db data
WORKDIR /ebrelayer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install ./cmd/ebrelayer

#
# Main
#
FROM alpine

# Copy the compiled binaires over.
COPY --from=build /go/bin/ebrelayer /usr/bin/ebrelayer
