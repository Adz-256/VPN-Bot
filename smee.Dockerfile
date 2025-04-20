FROM node:20

RUN npm install -g npm@latest

RUN npm install -g smee-client

CMD smee -t http://bot:3000 -u https://smee.io/xGrAI9x5GjgmIsl
