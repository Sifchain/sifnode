#!/bin/bash

[[ ! -z "$(docker ps -qaf name=sif-ui-stack)" ]] && docker stop sif-ui-stack && docker rm -f sif-ui-stack