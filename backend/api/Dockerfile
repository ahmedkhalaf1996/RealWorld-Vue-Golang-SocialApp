FROM golang:1.23.0 AS builder 
WORKDIR /usr/src/app 
COPY backend/api/go.mod ./
COPY backend/api/go.sum ./
RUN go mod download 
COPY backend/api . 
RUN CGO_ENABLED=0 GDOS=linux go build -o main ./cmd

FROM alpine:latest 
RUN apk --no-cache add ca-certificates 
WORKDIR /root/
COPY --from=builder /usr/src/app/main . 
COPY --from=builder /usr/src/app/.env . 
EXPOSE 5000
ENTRYPOINT [ "./main" ]






