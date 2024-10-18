# 1: build stage
FROM node:23-bookworm as build 

WORKDIR /app

COPY ./package*.json .

RUN npm ci

COPY . .

RUN npm run build

# 2: final stage 
FROM nginx:1.27.2-bookworm 

COPY --from=build /app/dist /usr/share/nginx/html


