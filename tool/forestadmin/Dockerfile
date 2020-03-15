FROM node:10-alpine
WORKDIR /usr/src/app
RUN npm install lumber-cli -g -s
COPY package*.json ./
RUN npm install -s
COPY . .
EXPOSE 3310
CMD ["npm", "start"]
