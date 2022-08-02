# Build Stage
FROM golang:1.18-alpine3.16 AS builder

RUN apk update && apk add --no-cache git && apk add curl

WORKDIR /app 

COPY go.mod go.sum ./

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

COPY . .

RUN go build -o main main.go 

# Run Stage 
FROM alpine:3.16

RUN addgroup app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/migrate.linux-amd64 ./migrate

COPY app.env start.sh wait-for.sh ./

COPY db/migration ./migration

EXPOSE 8080

CMD [ "/main" ]

ENTRYPOINT [ "/start.sh" ]
