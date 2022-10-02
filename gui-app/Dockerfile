FROM  node:16.4.2-alpine3.14 as node-base
WORKDIR /usr/app
COPY ./gui-app ./

RUN yarn

RUN npm build

CMD ["next", "start"]
