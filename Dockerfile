FROM golang:1.25-alpine

RUN go install github.com/air-verse/air@v1.63.4

WORKDIR /app

CMD ["air", "-c", "air.toml", "--", "serve"]