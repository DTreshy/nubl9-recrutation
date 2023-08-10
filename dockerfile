FROM golang:1.21.0-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o ./random ./cmd

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/random .
COPY --from=builder /app/build/config.yml .

RUN apk update && apk upgrade
RUN rm -rf /var/cashe/apk/*

EXPOSE 8080

CMD ["./random"]