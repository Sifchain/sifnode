#!/bin/bash

docker ps -q -f name=sifdevstack01 && docker stop sifdevstack01 && docker rm sifdevstack01

# this runs a docker image built by the build command
# the image temporarily will be pulled from dockerhub 
# but at a point soon will be transferred to our GH actions repo
docker run \
  -p 1317:1317 \
  -p 7545:7545 \
  -p 26656:26656 \
  -p 26657:26657 \
  --name sifdevstack01 \
  --platform linux/amd64 \
  ryardley/sifdevstack01:experimental