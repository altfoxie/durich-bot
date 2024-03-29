FROM golang:alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o durich-bot

FROM alpine

WORKDIR /app

COPY --from=builder /src/durich-bot .

ENV DB="./data/bot.db"
ENV SESSION="./data/session.json"

ENTRYPOINT ["/app/durich-bot"]