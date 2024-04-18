# Build stage
FROM golang:1.21.1 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Final stage
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=build /app/main .
CMD ["./main"]