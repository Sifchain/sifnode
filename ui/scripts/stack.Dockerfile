#
# GO Build
#
FROM golang:1.15 AS build-go

ARG chainnet

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
# Runtime
#
FROM node:14-alpine

EXPOSE 1317
EXPOSE 7545
EXPOSE 5000
EXPOSE 26656
EXPOSE 26657

RUN apk update && apk add curl jq bash

# Copy the compiled binaires over.
COPY --from=build-go /go/bin/ebrelayer /usr/bin/ebrelayer
COPY --from=build-go /go/bin/sifnoded /usr/bin/sifnoded
COPY --from=build-go /go/bin/sifgen /usr/bin/sifgen

# Required for ebrelayer
COPY --from=build-go /sif/cmd/ebrelayer /sif/cmd/ebrelayer

WORKDIR /sif/ui

COPY ./ui/package.json ./package.json
COPY ./ui/yarn.lock ./yarn.lock
COPY ./ui/chains/eth/package.json ./chains/eth/package.json
COPY ./ui/chains/eth/yarn.lock ./chains/eth/yarn.lock
COPY ./smart-contracts/package.json ../smart-contracts/package.json
COPY ./smart-contracts/yarn.lock ../smart-contracts/yarn.lock

RUN yarn install --frozen-lockfile --silent
RUN cd ./chains/eth && yarn install --frozen-lockfile --silent
RUN cd ../smart-contracts && yarn install --frozen-lockfile --silent

COPY ./ui/chains ./chains
COPY ./smart-contracts ../smart-contracts
COPY ./test/test-tables ../test/test-tables
COPY ./ui/scripts ./scripts

RUN ./scripts/build.sh

CMD ./scripts/start.sh