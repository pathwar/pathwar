FROM node:10
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install lumber-cli -g -s
RUN npm install -s
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
