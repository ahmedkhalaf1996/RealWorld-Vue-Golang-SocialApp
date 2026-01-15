FROM node:20.18.0-alpine AS build 
WORKDIR /usr/src/app 

COPY package.json  package-lock.json ./ 

RUN npm install 
COPY . . 
RUN  npm run build 

FROM nginx:1.15.8-alpine 
COPY --from=build /usr/src/app/dist /usr/share/nginx/html  
COPY .nginx/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD [ "nginx" ,"-g" , "daemon off;" ]
