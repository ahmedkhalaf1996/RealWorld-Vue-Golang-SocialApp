FROM golang:1.23.2 AS builder 
WORKDIR /usr/src/app 
COPY backend/realTimeChat/go.mod ./
COPY backend/realTimeChat/go.sum ./
RUN go mod download 
COPY backend/realTimeChat . 
RUN CGO_ENABLED=0 GDOS=linux go build -o main . 

FROM alpine:latest 
RUN apk --no-cache add ca-certificates 
WORKDIR /root/
COPY --from=builder /usr/src/app/main . 
COPY --from=builder /usr/src/app/.env . 
EXPOSE 8001
CMD [ "./main" ]
