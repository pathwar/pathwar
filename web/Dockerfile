FROM node:11-stretch

WORKDIR /app

COPY package*.json ./

RUN npm install
#RUN npm ci --only=production

COPY . .

EXPOSE 3000 3001
CMD [ "npm", "start" ]

# FIXME: add build step for production
