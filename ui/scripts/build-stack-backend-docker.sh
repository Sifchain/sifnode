#!/bin/bash

set -e


COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)

IMAGE_NAME=ghcr.io/sifchain/sifnode/ui-stack:$COMMIT
STABLE_TAG=ghcr.io/sifchain/sifnode/ui-stack:${BRANCH/\//__}

./scripts/ensure-docker-logged-in.sh

# reverse grep for go.mod because on CI this can be altered by installing go dependencies
if [[ -z "$CI" && ! -z "$(git status --porcelain --untracked-files=no)" ]]; then 
  echo "Git workspace must be clean to save git commit hash"
  exit 1
fi

echo "Github Registry Login found."
echo "Building new container..."

ROOT=$(pwd)/..

echo "New image name: $IMAGE_NAME"

# Using buildkit to take advantage of local dockerignore files
export DOCKER_BUILDKIT=1

cd $ROOT && docker build -f ./ui/scripts/stack.Dockerfile -t $IMAGE_NAME .

if [[ ! -z "$CI" ]]; then
  docker tag $IMAGE_NAME $STABLE_TAG
fi

docker push $IMAGE_NAME


# echo $IMAGE_NAME > $ROOT/ui/scripts/latest

# echo "Commit the ./ui/scripts/latest file to git"