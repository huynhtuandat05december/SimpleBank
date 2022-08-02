# Build stage
FROM golang:1.18-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.15
WORKDIR /usr/src/app
COPY --from=builder /app/main .
COPY .env .
COPY db/migration ./db/migration
COPY wait-for-it.sh .
COPY docker-entrypoint.sh .
RUN apk update && apk add bash
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh


EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD [ "/usr/src/app/main" ]


