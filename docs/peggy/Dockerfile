#
# Peggy docs.
#
FROM node:alpine

WORKDIR /var/www

ADD genesis.json genesis.json
ADD package.json package.json
ADD docs docs

RUN yarn install

EXPOSE 3000
