# Build Stage
FROM golang:1.20-alpine AS builder

WORKDIR /app 

COPY go.mod go.sum ./
COPY . .

# Download golang-migrate binary to run db migration in the container before starting the app
# RUN apk update && apk add --no-cache git && apk add curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

RUN go build -o main main.go 

# Run Stage 
FROM alpine:3.18

RUN addgroup app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env start.sh wait-for.sh ./
COPY db/migration ./db/migration

# COPY --from=builder /app/migrate.linux-amd64 ./migrate

EXPOSE 8080

CMD [ "/main" ]
ENTRYPOINT [ "/start.sh" ]
