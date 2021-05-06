#!/bin/bash

# If CI don't use tty - this is so we can use Ctrl C to cancel the stack script localy
[ -z "$CI" ] && IT="-it"

# get latest image name from latest file
IMAGE_NAME=$(cat ./scripts/latest)

# kill all other docker processes
docker ps -q -f name=sif-ui-stack && docker stop sif-ui-stack && docker rm sif-ui-stack

# this runs a docker image built by the build command
# the image temporarily will be pulled from dockerhub 
# but at a point soon will be transferred to our GH actions repo
docker run $IT \
  -p 1317:1317 \
  -p 7545:7545 \
  -p 26656:26656 \
  -p 26657:26657 \
  --name sif-ui-stack \
  --platform linux/amd64 \
  $IMAGE_NAME