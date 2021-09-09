#
# Build
#
FROM golang:1.16 AS build

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

EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

COPY --from=build /go/bin/sifnoded /usr/bin/sifnoded
COPY --from=build /go/bin/sifgen /usr/bin/sifgen
