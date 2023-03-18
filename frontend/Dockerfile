FROM node:16.17-alpine3.15 as builder

WORKDIR /usr/src/app

# install dependencies and avoid `node-gyp rebuild` errors
COPY ./package.json .
# COPY ./package-lock.json .
COPY ./yarn.lock .

# Install package
RUN yarn

# Copy source code
COPY ./ .

# Build
RUN yarn build

FROM nginx:1.19.0-alpine

WORKDIR /usr/share/nginx/html

COPY --from=builder /usr/src/app/build .

RUN rm /etc/nginx/conf.d/default.conf
COPY ./nginx.conf /etc/nginx/conf.d

ENTRYPOINT ["nginx", "-g", "daemon off;"]
