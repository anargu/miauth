FROM node:12-alpine

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY package*.json ./

# RUN apk update \
#     && apk --no-cache --update add build-base \
#     && apk --no-cache --update add python

RUN npm install
# RUN apk --no-cache add --virtual builds-deps build-base python && npm install && apk del builds-deps
# RUN apk --no-cache add --virtual native-deps \
#   g++ gcc libgcc libstdc++ linux-headers autoconf automake make nasm python git && \
#   npm install --quiet node-gyp -g

# RUN npm rebuild bcrypt --build-from-source


# If you are building your code for production
# RUN npm ci --only=production

# Bundle app source
COPY . .

CMD [ "npm", "run", "start" ]