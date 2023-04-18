FROM golang:alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o durich-bot

FROM alpine

WORKDIR /app

COPY --from=builder /src/durich-bot .
RUN touch session.json

ENTRYPOINT ["/app/durich-bot"]