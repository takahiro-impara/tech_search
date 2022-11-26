FROM node:16.15.0-alpine
WORKDIR /usr/src/app
COPY . .
EXPOSE 3000

RUN npm run build

CMD ["npm", "start"]