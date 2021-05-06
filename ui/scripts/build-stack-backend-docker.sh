#!/bin/bash


# this is temporarily being published on rudis docker hub account - @ryarfley
# TODO: investigate publishing on sifchain repository (possibly GH repo)


NOW=$(date +%s)
IMAGE_NAME=ghcr.io/sifchain/sifnode/ui-stack:$NOW
cd .. && docker build -f ./ui/scripts/stack.Dockerfile -t $IMAGE_NAME .
docker push $IMAGE_NAME

echo $IMAGE_NAME > ./scripts/latest