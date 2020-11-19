FROM node:14.11.0 as builder

RUN mkdir -p /ui/app
RUN mkdir -p /ui/core

ARG SIFNODE_API
ENV VUE_APP_SIFNODE_API=$SIFNODE_API

COPY ./package.json /ui/package.json
COPY ./yarn.lock /ui/yarn.lock
COPY ./app/package.json /ui/app/package.json
COPY ./app/yarn.lock /ui/app/yarn.lock
COPY ./core/package.json /ui/core/package.json
COPY ./core/yarn.lock /ui/core/yarn.lock
WORKDIR /ui
RUN yarn

COPY . /ui
RUN yarn core:build
RUN yarn app:build

FROM flashspys/nginx-static
RUN apk update && apk upgrade
RUN rm -rf /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder  /ui/app/dist /static
