FROM node:16-alpine

WORKDIR /app

# Copy the package.json and package-lock.json files to the working directory
COPY package*.json ./

# Install dependencies
RUN yarn install

COPY . .

EXPOSE 3000

CMD ["yarn", "start"]