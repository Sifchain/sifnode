FROM node:14.11.0 as builder
COPY . /ui
WORKDIR /ui
RUN yarn
RUN yarn core:build
RUN yarn app:build

FROM flashspys/nginx-static
RUN apk update && apk upgrade
RUN rm -rf /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder  /ui/app/dist /static
