# Build stage
FROM --platform=linux/amd64 golang:1.21.1 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=build /app/main .
CMD ["./main"]