# Build stage
FROM --platform=linux/amd64 golang:1.21.1 AS build
WORKDIR /app
COPY . .
RUN go build -o main .

# Final stage
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=build /app/main .
RUN ls -l
CMD ["./main"]