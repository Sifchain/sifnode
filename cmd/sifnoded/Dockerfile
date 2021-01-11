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
RUN go get github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor

COPY . .

RUN make install

#
# Main
#
FROM alpine

ENV PACKAGES curl jq bind-tools expect

EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

RUN apk add --update --no-cache $PACKAGES

# Copy the compiled binaires over.
COPY --from=build /go/bin/cosmovisor /usr/bin/cosmovisor
COPY --from=build /go/bin/sifnoded /usr/bin/sifnoded
COPY --from=build /go/bin/sifnodecli /usr/bin/sifnodecli
COPY --from=build /go/bin/sifgen /usr/bin/sifgen

CMD ["sifnoded", "start", "--rpc.laddr", "tcp://0.0.0.0:26657"]
