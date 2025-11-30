FROM node:20.18.0-alpine AS build 
WORKDIR /usr/src/app 
ARG VUE_APP_APIUrl
ARG VUE_APP_RealTimeNotificationUrl
ARG VUE_APP_RealTimeChatUrl

# Set as environment variables for the build
ENV VUE_APP_APIUrl=$VUE_APP_APIUrl
ENV VUE_APP_RealTimeNotificationUrl=$VUE_APP_RealTimeNotificationUrl
ENV VUE_APP_RealTimeChatUrl=$VUE_APP_RealTimeChatUrl
COPY frontend/package.json  frontend/package-lock.json ./ 
RUN npm install 

COPY frontend . 
RUN  npm run build 
FROM nginx:1.15.8-alpine 
COPY --from=build /usr/src/app/dist /usr/share/nginx/html  
COPY frontend/.nginx/nginx.conf /etc/nginx/nginx.conf

