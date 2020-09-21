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
ENV PACKAGES curl jq bind-tools
FROM alpine

RUN apk add --update --no-cache $PACKAGES

# Copy the compiled binaires over.
COPY --from=build /go/bin/sifnoded /usr/bin/sifd
COPY --from=build /go/bin/sifnodecli /usr/bin/sifcli

CMD ["sifd", "start"]
