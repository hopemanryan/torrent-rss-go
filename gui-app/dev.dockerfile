FROM  node:16.4.2-alpine3.14 as node-base
WORKDIR /usr/app
COPY ./gui-app ./

RUN npm cache clean --force

RUN npm i

RUN npm i --save @nestjs/core @nestjs/common rxjs reflect-metadata

CMD ["npm", "run", "dev"]
