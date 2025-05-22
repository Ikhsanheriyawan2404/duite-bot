FROM node:20-alpine

WORKDIR /bot

COPY package*.json ./
RUN npm install

COPY . .

CMD ["node", "index.js"]
