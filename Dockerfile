FROM --platform=linux/amd64 golang:1.21.1 AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM --platform=linux/amd64 golang:1.21.1

WORKDIR /app

COPY --from=build /app/main .
COPY config.yml .

CMD ["./main"]