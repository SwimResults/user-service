# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY user-service /app/service
RUN chmod +x /app/service

ENV SR_USER_PORT=8080

EXPOSE 8080

ENTRYPOINT [ "./service" ]
