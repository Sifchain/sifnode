#!/bin/bash

docker run -it \
  -p 1317:1317 \
  -p 7545:7545 \
  -p 5000:5000 \
  -p 26656:26656 \
  -p 26657:26657 \
  stack:latest