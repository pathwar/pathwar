FROM node:10-alpine
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install lumber-cli -g -s
RUN npm install -s
COPY . .
EXPOSE 3310
CMD ["npm", "start"]
