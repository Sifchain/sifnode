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
FROM alpine

ENV PACKAGES curl jq bind-tools supervisor

EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

RUN apk add --update --no-cache $PACKAGES

# Copy the smart contract ABIs over.
ADD cmd/ebrelayer/contract/generated  /sif/cmd/ebrelayer/contract/generated

# Copy the compiled binaires over.
COPY --from=build /go/bin/sifnoded /usr/bin/sifnoded
# COPY --from=build /go/bin/ebrelayer /usr/bin/ebrelayer
