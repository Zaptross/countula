FROM node:16-alpine3.11

# create folder
WORKDIR /usr/src/countula

# copy package lock and install prodlike build
COPY ./package*.json ./
RUN npm ci --only=production

# copy built code into build
COPY ./dist/* ./

# set startup command
CMD ["node", "index.js"]